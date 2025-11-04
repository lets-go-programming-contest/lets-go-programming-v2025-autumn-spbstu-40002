package parser

import (
	"encoding/xml"
	"fmt"
	"os"
)

type CurrencyItem struct {
	Name string  `xml:"Name" json:"name"`
	Rate float64 `xml:"Rate" json:"rate"`
}

type CurrencyList struct {
	Items []CurrencyItem `xml:"Item"`
}

func ParseXML(path string) ([]CurrencyItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open xml file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	var raw CurrencyList
	if err := xml.NewDecoder(file).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode xml: %w", err)
	}

	return raw.Items, nil
}
