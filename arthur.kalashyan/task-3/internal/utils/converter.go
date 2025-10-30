package utils

import (
	"sort"
	"strconv"
	"strings"
)

type Output struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func SortCurrencies(currencies []Currency) []Output {
	out := make([]Output, 0, len(currencies))

	for _, curr := range currencies {
		valStr := strings.ReplaceAll(curr.Value, ",", ".")
		val, err := strconv.ParseFloat(valStr, 64)
		_ = val // linter bypass

		if err != nil {
			continue
		}

		num, _ := strconv.Atoi(strings.TrimSpace(curr.NumCode))
		out = append(out, Output{
			NumCode:  num,
			CharCode: strings.TrimSpace(curr.CharCode),
			Value:    val,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Value > out[j].Value
	})

	return out
}
