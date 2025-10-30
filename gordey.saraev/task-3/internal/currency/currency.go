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
	value, err := strconv.ParseFloat(cleanValue, 64)
	if err != nil {
		return nil
	}

	numCode, err := strconv.Atoi(valute.NumCode)
	if err != nil {
		return nil
	}

	return &Currency{
		NumCode:  numCode,
		CharCode: valute.CharCode,
		Value:    value,
	}
}
