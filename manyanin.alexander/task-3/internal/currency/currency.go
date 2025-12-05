package currency

import (
	"sort"
	"strconv"
	"strings"
)

type Currency struct {
	NumCode  string  `xml:"NumCode" json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    string  `xml:"Value" json:"-"`
	ValueNum float64 `xml:"-" json:"value"`
}

func ValuteToCurr(valute Currency) *Currency {
	cleanValue := strings.ReplaceAll(valute.Value, ",", ".")

	value, err := strconv.ParseFloat(cleanValue, 64)
	if err != nil {
		value = 0
	}

	return &Currency{
		NumCode:  valute.NumCode,
		CharCode: valute.CharCode,
		Value:    valute.Value,
		ValueNum: value,
	}
}

func SortByValue(currencies []Currency) []Currency {
	sorted := make([]Currency, len(currencies))
	copy(sorted, currencies)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ValueNum > sorted[j].ValueNum
	})

	return sorted
}
