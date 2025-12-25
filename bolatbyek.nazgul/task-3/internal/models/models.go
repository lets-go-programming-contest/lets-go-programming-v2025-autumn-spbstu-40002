package models

import (
	"encoding/xml"
	"sort"
	"strconv"
	"strings"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	XMLName  xml.Name `xml:"Valute"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Value    string   `xml:"Value"`
}

type CurrencyOutput struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func ConvertAndSort(valCurs *ValCurs) []CurrencyOutput {
	outputCurrencies := make([]CurrencyOutput, 0, len(valCurs.Valutes))

	for _, valute := range valCurs.Valutes {
		valueStr := strings.Replace(valute.Value, ",", ".", 1)

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			// В случае ошибки парсинга значения считаем его равным 0,
			// но не паникуем и не выбрасываем запись.
			value = 0
		}

		numCode, err := strconv.Atoi(valute.NumCode)
		if err != nil {
			// Аналогично для некорректного числового кода — используем 0 по умолчанию.
			numCode = 0
		}

		outputCurrencies = append(outputCurrencies, CurrencyOutput{
			NumCode:  numCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	sort.Slice(outputCurrencies, func(i, j int) bool {
		return outputCurrencies[i].Value > outputCurrencies[j].Value
	})

	return outputCurrencies
}
