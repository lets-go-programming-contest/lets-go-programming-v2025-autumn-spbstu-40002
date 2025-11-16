package utils

import (
	"sort"

	"github.com/slendycs/go-lab-3/internal/parsers"
)

func SortVal(data *parsers.ValStruct) {
	sort.Slice(data.Valute, func(i, j int) bool {
		return data.Valute[i].Value > data.Valute[j].Value
	})
}
