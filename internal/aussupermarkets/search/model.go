// Package search defines the in-store search contract and provider orchestration.
package search

import (
	"context"
	"errors"
	"time"
)

type Retailer string

const (
	Coles      Retailer = "coles"
	Woolworths Retailer = "woolworths"
	ALDI       Retailer = "aldi"
)

// Retailers returns a fresh slice in the output order, independent of flag order.
func Retailers() []Retailer { return []Retailer{Coles, Woolworths, ALDI} }

type Sort string

const (
	SortPrice Sort = "price"
	SortName  Sort = "name"
	InStore        = "in_store"
)

// Money is exact: one minor unit is one Australian cent. No floating-point prices.
type Money struct {
	MinorUnits int64  `json:"minor_units"`
	Currency   string `json:"currency"`
}

type UnitPrice struct {
	Amount Money  `json:"amount"`
	Unit   string `json:"unit"`
}

type Product struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Price        *Money     `json:"price"`       // null is unknown; zero is a real zero price
	PriceBasis   string     `json:"price_basis"` // pack, weight, member or multibuy
	PackSize     string     `json:"pack_size,omitempty"`
	UnitPrice    *UnitPrice `json:"unit_price,omitempty"`
	Conditions   string     `json:"conditions,omitempty"` // any nonempty value excludes standard-pack comparability
	Availability string     `json:"availability,omitempty"`
	URL          string     `json:"url,omitempty"`
	ObservedAt   *time.Time `json:"observed_at,omitempty"`
}

type Location struct {
	Channel  string `json:"channel"`
	Postcode string `json:"postcode,omitempty"`
	StoreID  string `json:"store_id,omitempty"`
}

type Source struct {
	Name        string     `json:"name"`
	URL         string     `json:"url,omitempty"`
	RetrievedAt *time.Time `json:"retrieved_at,omitempty"`
	ObservedAt  *time.Time `json:"observed_at,omitempty"`
}

type Coverage struct {
	// complete_catalogue, sample or unknown. A sorted sample is not global top-N.
	Scope     string `json:"scope"`
	Truncated bool   `json:"truncated"`
}

type Request struct {
	Query      string
	Location   Location
	Sort       Sort
	MaxResults int // per retailer; adapters must not relevance-limit before sorting
}

type Result struct {
	Products []Product
	Location Location
	Source   Source
	Coverage Coverage
}

// Providers must honour cancellation and use only permitted in-store sources.
// Products share the returned Location; split/ambiguous-store data is not valid.
// Return complete candidates or a truthful Coverage describing the sampled set.
type Provider interface {
	Retailer() Retailer
	Search(context.Context, Request) (Result, error)
}

var (
	ErrUnconfigured = errors.New("provider unconfigured")
	ErrUnavailable  = errors.New("provider unavailable")
)

type Options struct {
	Query      string
	MaxResults int
	Sort       Sort
	Postcode   string
	Stores     map[Retailer]string
	Timeout    time.Duration // per provider, not a whole-invocation budget
}

type Group struct {
	Retailer          Retailer  `json:"retailer"`
	Status            string    `json:"status"`
	Error             string    `json:"error,omitempty"`
	RequestedLocation Location  `json:"requested_location"`
	Location          Location  `json:"location"`
	Source            Source    `json:"source"`
	Coverage          Coverage  `json:"coverage"`
	CandidateCount    int       `json:"candidate_count"`
	Limited           bool      `json:"limited"`
	Products          []Product `json:"products"`
}

type Report struct {
	Query      string  `json:"query"`
	Sort       Sort    `json:"sort"`
	MaxResults int     `json:"max_results_per_retailer"`
	Groups     []Group `json:"groups"`
}

func (r Report) Failed() bool {
	for _, g := range r.Groups {
		if g.Status != "ok" && g.Status != "empty" {
			return true
		}
	}
	return false
}
