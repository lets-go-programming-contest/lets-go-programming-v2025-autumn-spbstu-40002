package processor

import (
	"sort"

	"github.com/XShaygaND/task-3/internal/parser"
)

func OrganizeByRate(currencies []parser.CurrencyItem) []parser.CurrencyItem {
	sorted := make([]parser.CurrencyItem, len(currencies))
	copy(sorted, currencies)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Rate > sorted[j].Rate
	})

	return sorted
}
