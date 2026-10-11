package providers_test

import (
	"context"
	"testing"
	"time"

	"github.com/denwilliams/cli/internal/aussupermarkets/providers"
	"github.com/denwilliams/cli/internal/aussupermarkets/search"
)

func TestProductionHasNoSyntheticLiveResults(t *testing.T) {
	r, err := search.Run(context.Background(), providers.All(), search.Options{Query: "milk", MaxResults: 10, Sort: search.SortPrice, Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Groups) != 3 || !r.Failed() {
		t.Fatal(r)
	}
	for i, g := range r.Groups {
		if g.Retailer != search.Retailers()[i] || g.Status != "unconfigured" || len(g.Products) != 0 || g.Location.StoreID != "" || g.Source.Name != "" || g.Source.ObservedAt != nil {
			t.Fatalf("fake live result: %+v", g)
		}
	}
}
