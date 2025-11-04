package parser

import (
	"encoding/xml"
	"os"
)

type CurrencyItem struct {
	Name string  `xml:"Name"`
	Rate float64 `xml:"Rate"`
}

type CurrencyList struct {
	Items []CurrencyItem `xml:"Item"`
}

func ParseXML(path string) ([]CurrencyItem, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var list CurrencyList
	err = xml.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}
