package utils

import (
	"sort"
	"strconv"
	"strings"
)

func convertValutes(root valCurs) ([]OutputCurrency, error) {
	out := make([]OutputCurrency, 0, len(root.Valutes))
	for _, v := range root.Valutes {
		valStr := strings.ReplaceAll(strings.TrimSpace(v.Value), ",", ".")
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			continue
		}
		num, _ := strconv.Atoi(strings.TrimSpace(v.NumCode))
		out = append(out, OutputCurrency{
			NumCode:  num,
			CharCode: strings.TrimSpace(v.CharCode),
			Value:    val,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Value > out[j].Value
	})
	return out, nil
}
