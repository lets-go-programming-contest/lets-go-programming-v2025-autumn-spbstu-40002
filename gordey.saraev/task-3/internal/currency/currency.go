package currency

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type Currency struct {
	XMLName  xml.Name `xml:"Valute"`
	NumCode  string   `xml:"NumCode"  json:"num_code"`
	CharCode string   `xml:"CharCode" json:"char_code"`
	Value    string   `xml:"Value"    json:"value"`
}

func (c *Currency) ToNumeric() (numCode int, value float64) {
	num, err := strconv.Atoi(c.NumCode)
	if err != nil {
		num = 0
	}
	numCode = num

	cleanValue := strings.ReplaceAll(c.Value, ",", ".")
	val, err := strconv.ParseFloat(cleanValue, 64)
	if err != nil {
		val = 0
	}
	value = val

	return
}
