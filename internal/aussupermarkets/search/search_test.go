package search_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/denwilliams/cli/internal/aussupermarkets/search"
)

type fake struct {
	id search.Retailer
	fn func(context.Context, search.Request) (search.Result, error)
}

func (f fake) Retailer() search.Retailer { return f.id }
func (f fake) Search(ctx context.Context, r search.Request) (search.Result, error) {
	return f.fn(ctx, r)
}
func options() search.Options {
	return search.Options{Query: "milk", Sort: search.SortPrice, MaxResults: 2, Timeout: time.Second}
}
func result(req search.Request, products ...search.Product) search.Result {
	loc := req.Location
	if loc.StoreID == "" {
		loc.StoreID = "test-store"
	}
	return search.Result{Products: products, Location: loc, Source: search.Source{Name: "synthetic test fixture", URL: "https://example.invalid/feed"}, Coverage: search.Coverage{Scope: "complete_catalogue"}}
}
func pack(id, name string, cents int64) search.Product {
	return search.Product{ID: id, Name: name, Price: &search.Money{MinorUnits: cents, Currency: "AUD"}, PriceBasis: "pack"}
}
func run(t *testing.T, ctx context.Context, ps []search.Provider, o search.Options) search.Report {
	t.Helper()
	r, err := search.Run(ctx, ps, o)
	if err != nil {
		t.Fatal(err)
	}
	return r
}
func ids(g search.Group) []string {
	out := []string{}
	for _, p := range g.Products {
		out = append(out, p.ID)
	}
	return out
}

func TestConcurrentStartsAndCanonicalGroups(t *testing.T) {
	started := make(chan search.Retailer, 3)
	release := make(chan struct{})
	unlock := sync.OnceFunc(func() { close(release) })
	defer unlock()
	ps := []search.Provider{}
	for _, id := range []search.Retailer{search.ALDI, search.Woolworths, search.Coles} {
		ps = append(ps, fake{id, func(ctx context.Context, req search.Request) (search.Result, error) {
			started <- id
			select {
			case <-release:
				return result(req), nil
			case <-ctx.Done():
				return search.Result{}, ctx.Err()
			}
		}})
	}
	done := make(chan search.Report, 1)
	go func() { r, _ := search.Run(context.Background(), ps, options()); done <- r }()
	watchdog := time.NewTimer(time.Second)
	defer watchdog.Stop()
	for range ps {
		select {
		case <-started:
		case <-watchdog.C:
			t.Fatal("providers did not all start before release")
		}
	}
	unlock()
	var r search.Report
	select {
	case r = <-done:
	case <-watchdog.C:
		t.Fatal("providers did not return")
	}
	for i, id := range search.Retailers() {
		if r.Groups[i].Retailer != id || r.Groups[i].Status != "empty" {
			t.Fatalf("groups: %+v", r.Groups)
		}
	}
}

func TestSortBeforeLimitAndNoMutation(t *testing.T) {
	products := []search.Product{pack("b", "Same", 1002), {ID: "unknown", Name: "Unknown"}, pack("a", "Same", 1002), pack("free", "Free", 0), pack("expensive", "Apple", 9999), {ID: "member", Name: "Conditional", Price: &search.Money{MinorUnits: 1, Currency: "AUD"}, PriceBasis: "member", Conditions: "member only"}}
	original := append([]search.Product(nil), products...)
	for _, tc := range []struct {
		order search.Sort
		limit int
		want  []string
	}{
		{search.SortPrice, 2, []string{"free", "a"}},
		{search.SortPrice, 10, []string{"free", "a", "b", "expensive", "member", "unknown"}},
		{search.SortName, 3, []string{"expensive", "member", "free"}},
	} {
		t.Run(fmt.Sprintf("%s/%d", tc.order, tc.limit), func(t *testing.T) {
			o := options()
			o.Sort = tc.order
			o.MaxResults = tc.limit
			r := run(t, context.Background(), []search.Provider{fake{search.Coles, func(_ context.Context, req search.Request) (search.Result, error) {
				res := result(req, products...)
				res.Coverage = search.Coverage{Scope: "sample", Truncated: true}
				return res, nil
			}}}, o)
			g := r.Groups[0]
			if !reflect.DeepEqual(ids(g), tc.want) {
				t.Fatalf("got %v want %v", ids(g), tc.want)
			}
			if g.CandidateCount != len(products) || g.Limited != (len(products) > tc.limit) || !g.Coverage.Truncated || g.Coverage.Scope != "sample" {
				t.Fatalf("coverage lost: %+v", g)
			}
			if !reflect.DeepEqual(products, original) {
				t.Fatal("provider slice mutated")
			}
		})
	}
}

func TestDeterministicNameAndPriceTies(t *testing.T) {
	o := options()
	o.MaxResults = 10
	products := []search.Product{pack("z", "milk", 2), pack("b", "Milk", 2), pack("a", "Milk", 2), pack("n", "apple", 2)}
	for _, order := range []search.Sort{search.SortName, search.SortPrice} {
		o.Sort = order
		r := run(t, context.Background(), []search.Provider{fake{search.ALDI, func(_ context.Context, req search.Request) (search.Result, error) {
			return result(req, products...), nil
		}}}, o)
		if want := []string{"n", "a", "b", "z"}; !reflect.DeepEqual(ids(r.Groups[0]), want) {
			t.Fatalf("ties: %v", ids(r.Groups[0]))
		}
	}
}

func TestPartialFailurePreservesSuccessAndContext(t *testing.T) {
	o := options()
	o.Postcode = "0800"
	o.Stores = map[search.Retailer]string{search.Coles: "supplied-store"}
	ps := []search.Provider{
		fake{search.Coles, func(_ context.Context, req search.Request) (search.Result, error) {
			if req.Location.StoreID != "supplied-store" || req.Location.Postcode != "0800" || req.Location.Channel != search.InStore || req.MaxResults != 2 || req.Sort != search.SortPrice || req.Query != "milk" {
				t.Errorf("request: %+v", req)
			}
			return result(req, pack("p", "Milk", 345)), nil
		}},
		fake{search.Woolworths, func(_ context.Context, req search.Request) (search.Result, error) {
			if req.Location.StoreID != "" {
				t.Error("invented store")
			}
			return search.Result{}, errors.New("feed failed")
		}},
	}
	r := run(t, context.Background(), ps, o)
	if !r.Failed() || len(r.Groups) != 2 || r.Groups[0].Status != "ok" || len(r.Groups[0].Products) != 1 || r.Groups[1].Status != "error" || len(r.Groups[1].Products) != 0 {
		t.Fatalf("report: %+v", r)
	}
	if r.Groups[0].Location.StoreID != "supplied-store" || r.Groups[1].RequestedLocation.Postcode != "0800" || r.Groups[1].Location.StoreID != "" {
		t.Fatal("context lost or fabricated")
	}
}

func TestProviderErrorStates(t *testing.T) {
	for _, tc := range []struct {
		err    error
		status string
	}{
		{search.ErrUnavailable, "unavailable"}, {search.ErrUnconfigured, "unconfigured"}, {context.Canceled, "cancelled"}, {context.DeadlineExceeded, "timeout"}, {errors.New("failure"), "error"},
	} {
		t.Run(tc.status, func(t *testing.T) {
			r := run(t, context.Background(), []search.Provider{fake{search.Coles, func(context.Context, search.Request) (search.Result, error) { return search.Result{}, tc.err }}}, options())
			if r.Groups[0].Status != tc.status || !r.Failed() {
				t.Fatalf("%+v", r)
			}
		})
	}
}

func TestIsolatedTimeout(t *testing.T) {
	o := options()
	o.Timeout = 20 * time.Millisecond
	cancelled := make(chan struct{})
	ps := []search.Provider{fake{search.Coles, func(ctx context.Context, _ search.Request) (search.Result, error) {
		<-ctx.Done()
		close(cancelled)
		return search.Result{}, ctx.Err()
	}}, fake{search.ALDI, func(_ context.Context, req search.Request) (search.Result, error) {
		return result(req, pack("ok", "Milk", 123)), nil
	}}}
	r := run(t, context.Background(), ps, o)
	if r.Groups[0].Status != "timeout" || r.Groups[1].Status != "ok" {
		t.Fatalf("%+v", r)
	}
	select {
	case <-cancelled:
	case <-time.After(time.Second):
		t.Fatal("provider not cancelled")
	}
}

func TestParentCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	started := make(chan struct{}, 2)
	stopped := make(chan struct{}, 2)
	ps := []search.Provider{}
	for _, id := range []search.Retailer{search.Coles, search.ALDI} {
		ps = append(ps, fake{id, func(ctx context.Context, _ search.Request) (search.Result, error) {
			started <- struct{}{}
			<-ctx.Done()
			stopped <- struct{}{}
			return search.Result{}, ctx.Err()
		}})
	}
	done := make(chan search.Report, 1)
	go func() { r, _ := search.Run(ctx, ps, options()); done <- r }()
	for range ps {
		select {
		case <-started:
		case <-time.After(time.Second):
			t.Fatal("provider not started")
		}
	}
	cancel()
	select {
	case r := <-done:
		for _, g := range r.Groups {
			if g.Status != "cancelled" {
				t.Fatalf("%+v", r)
			}
		}
	case <-time.After(time.Second):
		t.Fatal("run did not cancel")
	}
	for range ps {
		select {
		case <-stopped:
		case <-time.After(time.Second):
			t.Fatal("provider did not stop")
		}
	}
}

func TestFaultyAdapterCannotHoldCallerPastDeadline(t *testing.T) {
	release := make(chan struct{})
	stopped := make(chan struct{})
	unlock := sync.OnceFunc(func() { close(release) })
	defer unlock()
	o := options()
	o.Timeout = 20 * time.Millisecond
	p := fake{search.Coles, func(_ context.Context, req search.Request) (search.Result, error) {
		<-release
		defer close(stopped)
		return result(req), nil
	}}
	done := make(chan search.Report, 1)
	go func() { r, _ := search.Run(context.Background(), []search.Provider{p}, o); done <- r }()
	select {
	case r := <-done:
		if r.Groups[0].Status != "timeout" {
			t.Fatal(r)
		}
	case <-time.After(time.Second):
		t.Fatal("caller hung")
	}
	// Unblock the deliberately faulty test adapter; never leave a leaked test worker.
	unlock()
	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("late adapter blocked")
	}
}

func TestPreCancelledContextDoesNotStartProvider(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := run(t, ctx, []search.Provider{fake{search.Coles, func(context.Context, search.Request) (search.Result, error) {
		t.Error("started after cancellation")
		return search.Result{}, nil
	}}}, options())
	if r.Groups[0].Status != "cancelled" {
		t.Fatal(r)
	}
}

func TestInvalidResultFailsClosed(t *testing.T) {
	cases := map[string]func(*search.Result){
		"online":                 func(r *search.Result) { r.Location.Channel = "delivery" },
		"missing store":          func(r *search.Result) { r.Location.StoreID = "" },
		"wrong store":            func(r *search.Result) { r.Location.StoreID = "wrong" },
		"wrong postcode":         func(r *search.Result) { r.Location.Postcode = "2000" },
		"missing source":         func(r *search.Result) { r.Source.Name = "" },
		"invalid source URL":     func(r *search.Result) { r.Source.URL = "javascript:bad" },
		"credentials URL":        func(r *search.Result) { r.Source.URL = "https://user:pass@example.invalid/feed" },
		"contradictory coverage": func(r *search.Result) { r.Coverage.Truncated = true },
		"unknown coverage":       func(r *search.Result) { r.Coverage.Scope = "" },
		"negative":               func(r *search.Result) { r.Products[0].Price.MinorUnits = -1 },
		"currency":               func(r *search.Result) { r.Products[0].Price.Currency = "USD" },
		"basis":                  func(r *search.Result) { r.Products[0].PriceBasis = "" },
		"conditional":            func(r *search.Result) { r.Products[0].PriceBasis = "multibuy" },
		"duplicate ID":           func(r *search.Result) { r.Products = append(r.Products, r.Products[0]) },
		"empty name":             func(r *search.Result) { r.Products[0].Name = " " },
		"unit price": func(r *search.Result) {
			r.Products[0].UnitPrice = &search.UnitPrice{Amount: search.Money{MinorUnits: 20, Currency: "AUD"}}
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			o := options()
			o.Postcode = "0800"
			o.Stores = map[search.Retailer]string{search.Coles: "requested"}
			r := run(t, context.Background(), []search.Provider{fake{search.Coles, func(_ context.Context, req search.Request) (search.Result, error) {
				res := result(req, pack("p", "Milk", 5))
				mutate(&res)
				return res, nil
			}}}, o)
			if r.Groups[0].Status != "error" || len(r.Groups[0].Products) != 0 {
				t.Fatalf("accepted invalid result: %+v", r)
			}
		})
	}
}

func TestInvalidOptionsAndRegistry(t *testing.T) {
	for name, mutate := range map[string]func(*search.Options){
		"empty query": func(o *search.Options) { o.Query = " " }, "limit": func(o *search.Options) { o.MaxResults = 0 }, "sort": func(o *search.Options) { o.Sort = "descending" }, "timeout": func(o *search.Options) { o.Timeout = 0 }, "postcode": func(o *search.Options) { o.Postcode = "800" }, "store": func(o *search.Options) { o.Stores = map[search.Retailer]string{search.Coles: " "} },
	} {
		t.Run(name, func(t *testing.T) {
			o := options()
			mutate(&o)
			if search.Validate(o) == nil {
				t.Fatal("accepted invalid options")
			}
		})
	}
	p := fake{search.Coles, func(context.Context, search.Request) (search.Result, error) {
		t.Error("invalid invocation started provider")
		return search.Result{}, nil
	}}
	for _, ps := range [][]search.Provider{nil, {nil}, {p, p}, {fake{id: "other"}}} {
		if _, err := search.Run(context.Background(), ps, options()); err == nil {
			t.Fatal("invalid registry accepted")
		}
	}
	o := options()
	o.Stores = map[search.Retailer]string{search.ALDI: "disabled"}
	if _, err := search.Run(context.Background(), []search.Provider{p}, o); err == nil {
		t.Fatal("disabled store accepted")
	}
}

func TestEmptyResultsStillRequireActualContextAndSource(t *testing.T) {
	for _, mutate := range []func(*search.Result){
		func(r *search.Result) { r.Location.StoreID = "" },
		func(r *search.Result) { r.Source = search.Source{} },
	} {
		r := run(t, context.Background(), []search.Provider{fake{search.Coles, func(_ context.Context, req search.Request) (search.Result, error) {
			res := result(req)
			mutate(&res)
			return res, nil
		}}}, options())
		if r.Groups[0].Status != "error" {
			t.Fatalf("unresolved empty result accepted: %+v", r)
		}
	}
}

func TestConditionalPackCannotHideUnconditionalOffer(t *testing.T) {
	for _, conditions := range []string{"members only", "buy 2 packs", " ", "\u00a0"} {
		for _, order := range []search.Sort{search.SortPrice, search.SortName} {
			for _, limit := range []int{1, 2, 3} {
				for _, reversed := range []bool{false, true} {
					t.Run(fmt.Sprintf("conditions=%q/%s/limit=%d/reversed=%t", conditions, order, limit, reversed), func(t *testing.T) {
						conditional := pack("conditional", "A conditional", 1)
						conditional.Conditions = conditions
						standard := pack("standard", "Standard", 500)
						unknown := search.Product{ID: "unknown", Name: "Z unknown"}
						products := []search.Product{conditional, unknown, standard}
						if reversed {
							products[0], products[2] = products[2], products[0]
						}
						original := append([]search.Product(nil), products...)
						o := options()
						o.Sort = order
						o.MaxResults = limit
						r := run(t, context.Background(), []search.Provider{fake{search.Coles, func(_ context.Context, req search.Request) (search.Result, error) {
							return result(req, products...), nil
						}}}, o)
						want := []string{"standard", "conditional", "unknown"}
						if order == search.SortName {
							want = []string{"conditional", "standard", "unknown"}
						}
						g := r.Groups[0]
						if g.Status != "ok" || !reflect.DeepEqual(ids(g), want[:limit]) {
							t.Fatalf("conditional pack displaced standard offer: got status=%s IDs=%v want=%v", g.Status, ids(g), want[:limit])
						}
						if g.CandidateCount != 3 || g.Limited != (limit < 3) {
							t.Fatalf("sort-before-limit metadata: %+v", g)
						}
						for _, p := range g.Products {
							if p.ID == "conditional" && (p.Conditions != conditions || p.Price.MinorUnits != 1) {
								t.Fatal("conditional offer details lost")
							}
						}
						if !reflect.DeepEqual(products, original) {
							t.Fatal("provider-owned candidates mutated")
						}
					})
				}
			}
		}
	}
}
