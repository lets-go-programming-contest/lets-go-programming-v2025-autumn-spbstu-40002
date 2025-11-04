// internal/processor/sorter.go
package processor

import (
	"sort"

	"github.com/XShaygaND/task-3/internal/types"
)

type CurrencySorter struct {
	Items []types.ProcessedCurrency
}

func (cs *CurrencySorter) SortByRateDescending() {
	sort.Slice(cs.Items, func(i, j int) bool {
		return cs.Items[i].Rate > cs.Items[j].Rate
	})
}

func OrganizeByRate(currencies []types.ProcessedCurrency) []types.ProcessedCurrency {
	sorter := CurrencySorter{Items: currencies}
	sorter.SortByRateDescending()

	return sorter.Items
}
