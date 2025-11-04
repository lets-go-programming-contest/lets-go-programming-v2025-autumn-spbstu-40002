package processor

import "sort"

type CurrencySorter struct {
	Items []ProcessedCurrency
}

func (cs *CurrencySorter) SortByRateDescending() {
	sort.Slice(cs.Items, func(i, j int) bool {
		return cs.Items[i].Rate > cs.Items[j].Rate
	})
}

func OrganizeByRate(currencies []ProcessedCurrency) []ProcessedCurrency {
	sorter := CurrencySorter{Items: currencies}
	sorter.SortByRateDescending()
	return sorter.Items
}
