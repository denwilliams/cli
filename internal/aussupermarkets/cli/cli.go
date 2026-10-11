// Package cli validates arguments and owns output/exit policy. No network calls
// or filesystem configuration are performed by the foundation command.
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/denwilliams/cli/internal/aussupermarkets/render"
	"github.com/denwilliams/cli/internal/aussupermarkets/search"
)

type storesFlag map[search.Retailer]string

func (s storesFlag) String() string { return "" }
func (s storesFlag) Set(value string) error {
	id, store, ok := strings.Cut(value, "=")
	if !ok || strings.TrimSpace(store) == "" || strings.TrimSpace(store) != store {
		return errors.New("store must be retailer=store-id")
	}
	r := search.Retailer(id)
	known := false
	for _, candidate := range search.Retailers() {
		if r == candidate {
			known = true
		}
	}
	if !known {
		return fmt.Errorf("unknown store retailer %q", id)
	}
	if _, exists := s[r]; exists {
		return fmt.Errorf("duplicate store for %s", r)
	}
	s[r] = store
	return nil
}

// Run returns 0 for help/all successful groups, 1 for provider/output failures,
// and 2 for invalid arguments. Failed groups remain in both output formats.
func Run(ctx context.Context, args []string, out, errOut io.Writer, registry []search.Provider) int {
	fs := flag.NewFlagSet("aussupermarkets", flag.ContinueOnError)
	fs.SetOutput(errOut)
	retailers := fs.String("providers", "coles,woolworths,aldi", "comma-separated enabled retailers (coles,woolworths,aldi)")
	limit := fs.Int("max-results", 10, "maximum results per supermarket (positive)")
	sortBy := fs.String("sort", "price", "ascending within-group sort: price or name")
	format := fs.String("format", "human", "output format: human or json")
	postcode := fs.String("postcode", "", "four-digit runtime postcode; no default location")
	timeout := fs.Duration("timeout", 10*time.Second, "positive per-provider timeout (e.g. 10s)")
	stores := storesFlag{}
	fs.Var(stores, "store", "runtime retailer=store-id; repeat for different retailers")
	fs.Usage = func() {
		fmt.Fprintln(errOut, "Usage: aussupermarkets [flags] query\n\nFoundation only: no live in-store feeds configured. Never substitutes online prices.\nQuote multiword queries. Flags must precede the query; use -- before a query beginning with '-'.")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	invalid := func(err error) int { fmt.Fprintln(errOut, "aussupermarkets:", err); return 2 }
	if fs.NArg() != 1 {
		return invalid(errors.New("supply one query argument (quote multiword queries); flags must precede it"))
	}
	if *format != "human" && *format != "json" {
		return invalid(errors.New("format must be human or json"))
	}
	opts := search.Options{Query: fs.Arg(0), MaxResults: *limit, Sort: search.Sort(*sortBy), Postcode: *postcode, Stores: stores, Timeout: *timeout}
	if err := search.Validate(opts); err != nil {
		return invalid(err)
	}
	registered := make(map[search.Retailer]search.Provider)
	for _, p := range registry {
		if p == nil {
			return invalid(errors.New("nil registered provider"))
		}
		id := p.Retailer()
		if _, ok := registered[id]; ok {
			return invalid(errors.New("duplicate registered provider"))
		}
		registered[id] = p
	}
	selected := []search.Provider{}
	seen := map[search.Retailer]bool{}
	for _, raw := range strings.Split(*retailers, ",") {
		id := search.Retailer(strings.TrimSpace(raw))
		p, ok := registered[id]
		if !ok {
			return invalid(fmt.Errorf("unknown or unavailable provider selection %q", raw))
		}
		if seen[id] {
			return invalid(fmt.Errorf("duplicate provider %q", id))
		}
		seen[id] = true
		selected = append(selected, p)
	}
	report, err := search.Run(ctx, selected, opts)
	if err != nil {
		return invalid(err)
	}
	if *format == "json" {
		err = render.JSON(out, report)
	} else {
		err = render.Human(out, report)
	}
	if err != nil {
		fmt.Fprintln(errOut, "aussupermarkets: output:", err)
		return 1
	}
	if report.Failed() {
		return 1
	}
	return 0
}
