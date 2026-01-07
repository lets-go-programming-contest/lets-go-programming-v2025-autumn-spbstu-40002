package utils

import (
	"sort"

	"github.com/braab/lets-go-programming-v2025-autumn-spbstu-40002/abdal.rahman/task-3/internal/parser"
)

func SortVal(data []parser.JSONData) {
	sort.Slice(data, func(i, j int) bool {
		return data[i].Value > data[j].Value
	})
}
