package utils

import (
	"strconv"
	"strings"

	"danila.onitshuk/task-3/internal/parser"
)

func parseValue(numberWithDot string) float64 {
	numberWithDot = strings.ReplaceAll(numberWithDot, ",", ".")

	num, err := strconv.ParseFloat(numberWithDot, 64)
	if err != nil {
		panic(err)
	}

	return num
}

func ToCurrency(data []parser.Valute) []parser.JSONData {
	result := make([]parser.JSONData, 0, len(data))

	for _, currency := range data {
		value := parseValue(currency.Value)
		result = append(result, parser.JSONData{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    value,
		})
	}

	return result
}
