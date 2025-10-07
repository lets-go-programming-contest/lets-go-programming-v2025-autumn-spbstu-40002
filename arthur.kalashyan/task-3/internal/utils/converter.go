package utils

import (
	"sort"
	"strconv"
	"strings"
)

type OutputItem struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func ConvertAndSort(values []cbrValue) []OutputItem {
	result := make([]OutputItem, 0, len(values))

	for _, val := range values {
		numCodeStr := strings.TrimSpace(val.NumCode)
		nominalStr := strings.TrimSpace(val.Nominal)
		valueStr := strings.TrimSpace(val.Value)

		if numCodeStr == "" || nominalStr == "" || valueStr == "" {
			continue
		}

		numCode, err := strconv.Atoi(numCodeStr)
		if err != nil {
			continue
		}

		nominal, err := strconv.Atoi(nominalStr)
		if err != nil || nominal == 0 {
			continue
		}

		valueStr = strings.ReplaceAll(valueStr, " ", "")
		valueStr = strings.ReplaceAll(valueStr, ",", ".")
		valueF, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}

		valuePerOne := valueF / float64(nominal)

		result = append(result, OutputItem{
			NumCode:  numCode,
			CharCode: strings.TrimSpace(val.CharCode),
			Value:    valuePerOne,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Value > result[j].Value
	})

	return result
}
