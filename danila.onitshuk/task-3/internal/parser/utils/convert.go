package utils

import (
	"strconv"
	"strings"

	"danila.onitshuk/task-3/internal/parser"
)

func parseValue(numberWithDot string) (float64, error) {
	numberWithDot = strings.ReplaceAll(numberWithDot, ",", ".")

	return strconv.ParseFloat(numberWithDot, 64)
}

func ToCurrency(data []parser.Valute) []parser.JsonData {
	result := make([]parser.JsonData, 0, len(data))

	for _, currency := range data {
		value, err := parseValue(currency.Value)
		if err != nil {
			panic(err)
		}

		result = append(result, parser.JsonData{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    value,
		})
	}

	return result
}
