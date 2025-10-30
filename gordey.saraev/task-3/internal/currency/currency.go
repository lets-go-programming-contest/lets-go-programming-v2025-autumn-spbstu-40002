package currency

import (
	"encoding/xml"
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

func ValuteToCurrency(valute Valute) *Currency {
	cleanValue := strings.Replace(valute.Value, ",", ".", -1)
	cleanValue = strings.TrimSpace(cleanValue)
	value, err := strconv.ParseFloat(cleanValue, 64)
	if err != nil {
		return nil
	}

	numCode := 0
	if valute.NumCode != "" {
		if parsed, err := strconv.Atoi(strings.TrimSpace(valute.NumCode)); err == nil {
			numCode = parsed
		}
	}

	if valute.CharCode == "" {
		return nil
	}

	return &Currency{
		NumCode:  numCode,
		CharCode: valute.CharCode,
		Value:    value,
	}
}
