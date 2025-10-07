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
	var result []OutputItem

	for _, val := range values {
		numCode, err := strconv.Atoi(strings.TrimSpace(val.NumCode))
		if err != nil {
			panic(err)
		}

		nominal, err := strconv.Atoi(strings.TrimSpace(val.Nominal))
		if err != nil {
			panic(err)
		}

		raw := strings.TrimSpace(val.Value)
		raw = strings.ReplaceAll(raw, " ", "")
		raw = strings.ReplaceAll(raw, ",", ".")
		valueF, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			panic(err)
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
