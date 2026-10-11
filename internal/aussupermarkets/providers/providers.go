// Package providers registers only genuine production provider states. There are
// deliberately no synthetic products or online-price fallbacks in this package.
package providers

import (
	"context"
	"fmt"

	"github.com/denwilliams/cli/internal/aussupermarkets/search"
)

type unconfigured struct{ id search.Retailer }

func (p unconfigured) Retailer() search.Retailer { return p.id }
func (p unconfigured) Search(ctx context.Context, _ search.Request) (search.Result, error) {
	if err := ctx.Err(); err != nil {
		return search.Result{}, err
	}
	return search.Result{}, fmt.Errorf("%w: %s has no verified permitted live in-store feed; no online prices substituted", search.ErrUnconfigured, p.id)
}

func All() []search.Provider {
	result := make([]search.Provider, 0, 3)
	for _, id := range search.Retailers() {
		result = append(result, unconfigured{id: id})
	}
	return result
}
