// Package render produces grouped human output and the stable JSON schema.
package render

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/denwilliams/cli/internal/aussupermarkets/search"
)

func JSON(w io.Writer, r search.Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func Money(m search.Money) string {
	return fmt.Sprintf("%s %d.%02d", m.Currency, m.MinorUnits/100, m.MinorUnits%100)
}

// Quote provider-controlled strings so terminal escapes/newlines cannot alter output.
func text(s string) string { return fmt.Sprintf("%q", s) }
func stamp(t *time.Time) string {
	if t == nil {
		return "unknown"
	}
	return t.UTC().Format(time.RFC3339)
}
func location(l search.Location) string {
	postcode, store := l.Postcode, l.StoreID
	if postcode == "" {
		postcode = "unspecified"
	}
	if store == "" {
		store = "unspecified"
	}
	return fmt.Sprintf("channel=%s postcode=%s store=%s", text(l.Channel), text(postcode), text(store))
}

func Human(w io.Writer, r search.Report) error {
	var b strings.Builder
	fmt.Fprintf(&b, "In-store search %s (ascending %s; max %d per retailer)\n", text(r.Query), r.Sort, r.MaxResults)
	b.WriteString("Pack prices are not quantity-normalized. In price sort, weight/conditional/unknown prices follow standard pack prices.\n")
	for _, g := range r.Groups {
		fmt.Fprintf(&b, "\n%s [%s]\n  Requested: %s\n", g.Retailer, g.Status, location(g.RequestedLocation))
		if g.Error != "" {
			fmt.Fprintf(&b, "  %s\n", text(g.Error))
			continue
		}
		fmt.Fprintf(&b, "  Actual: %s\n  Source: %s %s\n  Freshness: observed=%s retrieved=%s\n  Coverage: %s; provider-truncated=%t; candidates=%d; limited=%t\n", location(g.Location), text(g.Source.Name), text(g.Source.URL), stamp(g.Source.ObservedAt), stamp(g.Source.RetrievedAt), text(g.Coverage.Scope), g.Coverage.Truncated, g.CandidateCount, g.Limited)
		if len(g.Products) == 0 {
			b.WriteString("  No matches.\n")
		}
		for _, p := range g.Products {
			price := "price unknown"
			if p.Price != nil {
				price = Money(*p.Price) + " (" + p.PriceBasis + ")"
			}
			fmt.Fprintf(&b, "  - %s — %s", text(p.Name), price)
			if p.PackSize != "" {
				fmt.Fprintf(&b, "; pack=%s", text(p.PackSize))
			}
			if p.UnitPrice != nil {
				fmt.Fprintf(&b, "; unit=%s / %s", Money(p.UnitPrice.Amount), text(p.UnitPrice.Unit))
			}
			if p.Conditions != "" {
				fmt.Fprintf(&b, "; conditions=%s", text(p.Conditions))
			}
			if p.Availability != "" {
				fmt.Fprintf(&b, "; availability=%s", text(p.Availability))
			}
			fmt.Fprintf(&b, "; observed=%s", stamp(p.ObservedAt))
			if p.URL != "" {
				fmt.Fprintf(&b, "; url=%s", text(p.URL))
			}
			b.WriteByte('\n')
		}
	}
	_, err := io.WriteString(w, b.String())
	return err
}
