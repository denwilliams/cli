package search

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func validRetailer(id Retailer) bool {
	for _, r := range Retailers() {
		if r == id {
			return true
		}
	}
	return false
}

func ValidPostcode(p string) bool {
	if len(p) != 4 {
		return false
	}
	for _, c := range p {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func Validate(o Options) error {
	if strings.TrimSpace(o.Query) == "" {
		return errors.New("query must not be empty")
	}
	if o.MaxResults <= 0 {
		return errors.New("max-results must be positive")
	}
	if o.Sort != SortPrice && o.Sort != SortName {
		return errors.New("sort must be price or name")
	}
	if o.Timeout <= 0 {
		return errors.New("timeout must be positive")
	}
	if o.Postcode != "" && !ValidPostcode(o.Postcode) {
		return errors.New("postcode must contain exactly four digits")
	}
	for id, store := range o.Stores {
		if !validRetailer(id) || strings.TrimSpace(store) == "" || strings.TrimSpace(store) != store {
			return errors.New("store must have a known retailer and nonempty, trimmed ID")
		}
	}
	return nil
}

// Run concurrently searches selected providers and assembles canonical retailer order.
// Its deadlines bound the caller even if a faulty adapter is slow to return; adapters
// still MUST honour their context to release their own resources.
func Run(ctx context.Context, providers []Provider, o Options) (Report, error) {
	if err := Validate(o); err != nil {
		return Report{}, err
	}
	byID := make(map[Retailer]Provider)
	for _, p := range providers {
		if p == nil {
			return Report{}, errors.New("nil provider")
		}
		id := p.Retailer()
		if !validRetailer(id) {
			return Report{}, fmt.Errorf("unknown retailer %q", id)
		}
		if _, exists := byID[id]; exists {
			return Report{}, fmt.Errorf("duplicate retailer %q", id)
		}
		byID[id] = p
	}
	if len(byID) == 0 {
		return Report{}, errors.New("at least one provider must be enabled")
	}
	for id := range o.Stores {
		if _, ok := byID[id]; !ok {
			return Report{}, fmt.Errorf("store supplied for disabled retailer %q", id)
		}
	}
	report := Report{Query: strings.TrimSpace(o.Query), Sort: o.Sort, MaxResults: o.MaxResults, Groups: make([]Group, 0, len(byID))}
	channels := make([]<-chan Group, 0, len(byID))
	for _, id := range Retailers() {
		p, ok := byID[id]
		if !ok {
			continue
		}
		req := Request{Query: report.Query, Sort: o.Sort, MaxResults: o.MaxResults, Location: Location{Channel: InStore, Postcode: o.Postcode, StoreID: o.Stores[id]}}
		ch := make(chan Group, 1)
		channels = append(channels, ch)
		go func() { ch <- runProvider(ctx, p, req, o) }()
	}
	for _, ch := range channels {
		report.Groups = append(report.Groups, <-ch)
	}
	return report, nil
}

type outcome struct {
	result Result
	err    error
}

func runProvider(parent context.Context, p Provider, req Request, o Options) Group {
	g := Group{Retailer: p.Retailer(), RequestedLocation: req.Location, Location: Location{Channel: InStore}, Coverage: Coverage{Scope: "unknown"}, Products: []Product{}}
	ctx, cancel := context.WithTimeout(parent, o.Timeout)
	defer cancel()
	// A buffered channel lets a late cooperative provider return without blocking.
	ch := make(chan outcome, 1)
	if err := ctx.Err(); err != nil {
		setError(&g, err)
		return g
	}
	go func() { result, err := p.Search(ctx, req); ch <- outcome{result, err} }()
	var out outcome
	select {
	case <-ctx.Done():
		setError(&g, ctx.Err())
		return g
	case out = <-ch:
	}
	if err := ctx.Err(); err != nil {
		setError(&g, err)
		return g
	}
	if out.err != nil {
		setError(&g, out.err)
		return g
	}
	if err := validateResult(out.result, req.Location); err != nil {
		setError(&g, fmt.Errorf("invalid in-store result: %w", err))
		return g
	}
	g.Location, g.Source, g.Coverage = out.result.Location, out.result.Source, out.result.Coverage
	// Copy the slice: sorting must not mutate provider-owned storage.
	g.Products = append(g.Products, out.result.Products...)
	sort.Slice(g.Products, func(i, j int) bool { return less(g.Products[i], g.Products[j], o.Sort) })
	g.CandidateCount = len(g.Products)
	g.Limited = len(g.Products) > o.MaxResults
	if g.Limited {
		g.Products = g.Products[:o.MaxResults]
	}
	g.Status = "ok"
	if g.CandidateCount == 0 {
		g.Status = "empty"
	}
	return g
}

func setError(g *Group, err error) {
	g.Error = err.Error()
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		g.Status = "timeout"
	case errors.Is(err, context.Canceled):
		g.Status = "cancelled"
	case errors.Is(err, ErrUnconfigured):
		g.Status = "unconfigured"
	case errors.Is(err, ErrUnavailable):
		g.Status = "unavailable"
	default:
		g.Status = "error"
	}
}

// A pack total is standard only when it has no conditions. Treat even whitespace
// as a nonempty condition conservatively; adapters should not erase restrictions.
func comparable(p Product) bool {
	return p.Price != nil && p.PriceBasis == "pack" && p.Conditions == ""
}
func less(a, b Product, order Sort) bool {
	if order == SortPrice {
		ac, bc := comparable(a), comparable(b)
		if ac != bc {
			return ac
		}
		if ac && a.Price.MinorUnits != b.Price.MinorUnits {
			return a.Price.MinorUnits < b.Price.MinorUnits
		}
	}
	an, bn := strings.ToLower(a.Name), strings.ToLower(b.Name)
	if an != bn {
		return an < bn
	}
	if a.Name != b.Name {
		return a.Name < b.Name
	}
	return a.ID < b.ID
}

func validURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "https" || u.Scheme == "http") && u.Host != "" && u.User == nil
}
func validMoney(m Money) bool { return m.Currency == "AUD" && m.MinorUnits >= 0 }

func validateResult(r Result, requested Location) error {
	if r.Location.Channel != InStore {
		return errors.New("source channel is not in_store")
	}
	if r.Location.Postcode != "" && !ValidPostcode(r.Location.Postcode) {
		return errors.New("invalid actual postcode")
	}
	if requested.StoreID != "" && r.Location.StoreID != requested.StoreID {
		return errors.New("store does not match requested context")
	}
	if requested.Postcode != "" && r.Location.Postcode != requested.Postcode {
		return errors.New("postcode does not match requested context")
	}
	switch r.Coverage.Scope {
	case "complete_catalogue", "sample", "unknown":
	default:
		return errors.New("coverage must be explicit")
	}
	if strings.TrimSpace(r.Location.StoreID) == "" || strings.TrimSpace(r.Source.Name) == "" || !validURL(r.Source.URL) {
		return errors.New("successful results require an actual store and source name/URL")
	}
	if r.Coverage.Scope == "complete_catalogue" && r.Coverage.Truncated {
		return errors.New("truncated coverage cannot be a complete catalogue")
	}
	seen := make(map[string]bool)
	for _, p := range r.Products {
		if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.Name) == "" || seen[p.ID] {
			return errors.New("product IDs must be unique and names/IDs nonempty")
		}
		seen[p.ID] = true
		if p.URL != "" && !validURL(p.URL) {
			return errors.New("invalid product URL")
		}
		if p.Price != nil {
			if !validMoney(*p.Price) {
				return errors.New("price must be nonnegative AUD minor units")
			}
			switch p.PriceBasis {
			case "pack", "weight", "member", "multibuy":
			default:
				return errors.New("priced product requires a known price basis")
			}
			if (p.PriceBasis == "member" || p.PriceBasis == "multibuy") && strings.TrimSpace(p.Conditions) == "" {
				return errors.New("conditional price requires conditions")
			}
		}
		if p.UnitPrice != nil && (!validMoney(p.UnitPrice.Amount) || strings.TrimSpace(p.UnitPrice.Unit) == "") {
			return errors.New("invalid unit price")
		}
	}
	return nil
}
