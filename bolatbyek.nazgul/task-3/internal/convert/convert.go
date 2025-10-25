package convert

import (
	"sort"
	"github.com/bolatbyek/task-3/internal/cbr"
)

// Converter handles currency data conversion
type Converter struct{}

// NewConverter creates a new converter
func NewConverter() *Converter {
	return &Converter{}
}

// SortByValue sorts currencies by value in descending order
func (c *Converter) SortByValue(currencies []cbr.Currency) []cbr.Currency {
	sorted := make([]cbr.Currency, len(currencies))
	copy(sorted, currencies)
	
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	
	return sorted
}