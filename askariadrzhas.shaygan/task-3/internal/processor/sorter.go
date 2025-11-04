package processor

import (
	"sort"

	"github.com/XShaygaND/task-3/internal/parser"
)

type CurrencySorter struct {
	Items []parser.ProcessedCurrency
}

func (cs *CurrencySorter) SortByRateDescending() {
	sort.Slice(cs.Items, func(i, j int) bool {
		return cs.Items[i].Rate > cs.Items[j].Rate
	})
}

func OrganizeByRate(currencies []parser.ProcessedCurrency) []parser.ProcessedCurrency {
	sorter := CurrencySorter{Items: currencies}
	sorter.SortByRateDescending()

	return sorter.Items
}
