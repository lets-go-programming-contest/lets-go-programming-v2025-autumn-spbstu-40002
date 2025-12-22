package processor

import (
	"strconv"
	"strings"

	"github.com/stepanov.alexander/task-3/pkg/loader"
)

type CurrencyRate struct {
	NumCode   int     `json:"num_code"`
	CharCode  string  `json:"char_code"`
	Value     float64 `json:"value"`
}

func ProcessXML(valCurs *loader.ValCurs) ([]CurrencyRate, error) {
	var rates []CurrencyRate

	for _, valute := range valCurs.Valutes {
		numCode, err := strconv.Atoi(valute.NumCode)
		if err != nil {
			continue
		}

		valueStr := strings.TrimSpace(valute.Value)
		valueStr = strings.ReplaceAll(valueStr, ",", ".")

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}

		rates = append(rates, CurrencyRate{
			NumCode:  numCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	return rates, nil
}
