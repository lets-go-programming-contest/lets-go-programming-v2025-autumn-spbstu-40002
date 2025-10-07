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

func ConvertAndSort(vals []cbrValue) []OutputItem {
	var out []OutputItem
	for _, v := range vals {
		numCode, err := strconv.Atoi(strings.TrimSpace(v.NumCode))
		if err != nil {
			panic(err)
		}
		nominal, err := strconv.Atoi(strings.TrimSpace(v.Nominal))
		if err != nil {
			panic(err)
		}
		raw := strings.TrimSpace(v.Value)
		raw = strings.ReplaceAll(raw, " ", "")
		raw = strings.ReplaceAll(raw, ",", ".")
		valueF, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			panic(err)
		}
		valuePerOne := valueF / float64(nominal)
		out = append(out, OutputItem{
			NumCode:  numCode,
			CharCode: strings.TrimSpace(v.CharCode),
			Value:    valuePerOne,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Value > out[j].Value
	})
	return out
}
