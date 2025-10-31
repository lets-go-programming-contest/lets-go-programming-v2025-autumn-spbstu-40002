package utils

import (
	"strconv"
	"strings"

	"github.com/rachguta/task-3/internal/currency"
)

func ParseCurrencyValue(val string) float64 {
	if val == "" {
		return 0.0
	}

	val = strings.Replace(val, ",", ".", 1)

	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0.0
	}

	return floatVal
}

func ParseValutesToJSON(valutes []currency.Valute) []currency.JSONValute {
	result := make([]currency.JSONValute, 0, len(valutes))

	for _, cur := range valutes {
		result = append(result, currency.JSONValute{
			NumCode:  cur.NumCode,
			CharCode: cur.CharCode,
			Value:    ParseCurrencyValue(cur.Value),
		})
	}

	return result
}
