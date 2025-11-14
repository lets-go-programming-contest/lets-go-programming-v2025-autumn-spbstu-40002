package currency

import (
	"encoding/xml"
	"sort"
	"strconv"
	"strings"
)

type Valute struct {
	XMLName  xml.Name `xml:"Valute"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Value    string   `xml:"Value"`
}

type Currency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func ValuteToCurr(valute Valute) *Currency {
	cleanValue := strings.ReplaceAll(valute.Value, ",", ".")

	value, err := strconv.ParseFloat(cleanValue, 64)
	if err != nil {
		value = 0
	}

	numCode, err := strconv.Atoi(valute.NumCode)
	if err != nil {
		numCode = 0
	}

	return &Currency{
		NumCode:  numCode,
		CharCode: valute.CharCode,
		Value:    value,
	}
}

func SortByValue(currencies []Currency) []Currency {
	sorted := make([]Currency, len(currencies))
	copy(sorted, currencies)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	return sorted
}
