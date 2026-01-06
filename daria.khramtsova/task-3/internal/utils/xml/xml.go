package xml

import (
	"strconv"
	"strings"
)

type Currency struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

func (c Currency) FloatValue() float64 {
	value := strings.ReplaceAll(c.Value, ",", ".")
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(err)
	}
	return f
}
