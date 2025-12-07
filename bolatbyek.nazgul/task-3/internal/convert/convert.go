package convert

import (
	"sort"
	"strconv"
	"strings"

	"github.com/Nazkaaa/task-3/internal/cbr"
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
			// В случае ошибки парсинга значения считаем его равным 0,
			// но не паникуем и не выбрасываем запись.
			value = 0
		}

		numCode, err := strconv.Atoi(v.NumCode)
		if err != nil {
			// Аналогично для некорректного числового кода — используем 0 по умолчанию.
			numCode = 0
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
