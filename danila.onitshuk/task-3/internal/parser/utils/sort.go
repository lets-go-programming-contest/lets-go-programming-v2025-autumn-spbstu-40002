package utils

import (
	"sort"

	"danila.onitshuk/task-3/internal/parser"
)

func SortVal(data []parser.JSONData) {
	sort.Slice(data, func(i, j int) bool {
		return data[i].Value > data[j].Value
	})
}
