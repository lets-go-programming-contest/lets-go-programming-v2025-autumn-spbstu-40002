package processor

import (
	"sort"

	"github.com/XShaygaND/task-3/internal/parser"
)

// SortData sorts items by ID
func SortData(items []parser.Item) []parser.Item {
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})
	return items
}
