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
	var out []Output
	for _, c := range currencies {
		valStr := strings.ReplaceAll(c.Value, ",", ".")
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			continue
		}
		num, _ := strconv.Atoi(strings.TrimSpace(c.NumCode))
		out = append(out, Output{
			NumCode:  num,
			CharCode: strings.TrimSpace(c.CharCode),
			Value:    val,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Value > out[j].Value
	})
	return out
}
