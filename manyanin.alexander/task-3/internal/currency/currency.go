package currency

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Currency struct {
	XMLName  xml.Name `xml:"Valute"`
	NumCode  int      `json:"num_code"  xml:"NumCode"`
	CharCode string   `json:"char_code" xml:"CharCode"`
	Value    float64  `json:"value"     xml:"Value"`
}

func (c *Currency) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type currencyXML struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	var aux currencyXML
	if err := d.DecodeElement(&aux, &start); err != nil {
		return err
	}

	numCode, err := strconv.Atoi(aux.NumCode)
	if err != nil {
		return fmt.Errorf("failed to parse NumCode: %w", err)
	}

	c.NumCode = numCode

	cleanValue := strings.ReplaceAll(aux.Value, ",", ".")

	value, err := strconv.ParseFloat(cleanValue, 64)
	if err != nil {
		return fmt.Errorf("failed to parse Value: %w", err)
	}

	c.Value = value

	c.CharCode = aux.CharCode

	c.XMLName = start.Name

	return nil
}

func SortByValue(currencies []Currency) []Currency {
	sorted := make([]Currency, len(currencies))
	copy(sorted, currencies)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	return sorted
}
