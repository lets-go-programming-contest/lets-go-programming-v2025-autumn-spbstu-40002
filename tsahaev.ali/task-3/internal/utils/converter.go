package utils

import (
	"strconv"
	"strings"

	"github.com/Tsahaev/task-3/internal/currency"
)

func ParseCurrencyValue(val string) float64 {
	if val == "" {
		return 0
	}
	val = strings.Replace(val, ",", ".", 1)
	res, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0
	}
	return res
}

func ParseValutesToJSON(valutes []currency.Valute) []currency.JSONValute {
	out := make([]currency.JSONValute, 0, len(valutes))
	for _, v := range valutes {
		out = append(out, currency.JSONValute{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    ParseCurrencyValue(v.Value),
		})
	}
	return out
}
