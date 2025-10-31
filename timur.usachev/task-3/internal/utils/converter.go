package utils

import (
	"sort"
	"strconv"
	"strings"
)

func convertValutes(root valCurs) []OutputCurrency {
	result := make([]OutputCurrency, 0, len(root.Valutes))

	for _, valute := range root.Valutes {
		valueFloatStr := strings.ReplaceAll(strings.TrimSpace(valute.Value), ",", ".")
		valueFloat, err := strconv.ParseFloat(valueFloatStr, 64)
		if err != nil {
			continue
		}

		numCode, _ := strconv.Atoi(strings.TrimSpace(valute.NumCode))

		result = append(result, OutputCurrency{
			NumCode:  numCode,
			CharCode: strings.TrimSpace(valute.CharCode),
			Value:    valueFloat,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Value > result[j].Value
	})

	return result
}
