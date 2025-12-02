package utils

import (
	"sort"
	"strconv"
	"strings"
)

func convertValutes(root valCurs) []OutputCurrency {
	out := make([]OutputCurrency, 0, len(root.Valutes))

	for _, valute := range root.Valutes {
		valStr := strings.ReplaceAll(strings.TrimSpace(valute.Value), ",", ".")

		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			continue
		}

		num, _ := strconv.Atoi(strings.TrimSpace(valute.NumCode))

		out = append(out, OutputCurrency{
			NumCode:  num,
			CharCode: strings.TrimSpace(valute.CharCode),
			Value:    val,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Value > out[j].Value
	})

	return out
}
