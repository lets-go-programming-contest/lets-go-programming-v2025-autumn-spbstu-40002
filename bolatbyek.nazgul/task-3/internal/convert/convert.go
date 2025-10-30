package convert

import (
	"sort"
	"strconv"
	"strings"

	"lets-go-programming-v2025-autumn-spbstu-40002/internal/cbr"
)

type CurrencyOutput struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func ConvertAndSort(valCurs *cbr.ValCurs) []CurrencyOutput {
	var outputCurrencies []CurrencyOutput

	for _, v := range valCurs.Valutes {
		valueStr := strings.Replace(v.Value, ",", ".", 1)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			panic("Error parsing currency value: " + err.Error())
		}

		numCode, err := strconv.Atoi(v.NumCode)
		if err != nil {
			panic("Error parsing numeric code: " + err.Error())
		}

		outputCurrencies = append(outputCurrencies, CurrencyOutput{
			NumCode:  numCode,
			CharCode: v.CharCode,
			Value:    value,
		})
	}

	sort.Slice(outputCurrencies, func(i, j int) bool {
		return outputCurrencies[i].Value > outputCurrencies[j].Value
	})

	return outputCurrencies
}