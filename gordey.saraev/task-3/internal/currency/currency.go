package currency

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type Currency struct {
	XMLName  xml.Name `xml:"Valute"   json:"-"`
	NumCode  int      `xml:"NumCode"  json:"num_code"`
	CharCode string   `xml:"CharCode" json:"char_code"`
	Value    float64  `xml:"Value"    json:"value"`
}

func (c *Currency) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var raw struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}
	if err := decoder.DecodeElement(&raw, &start); err != nil {
		return fmt.Errorf("decode XML element: %w", err)
	}

	num, err := strconv.Atoi(raw.NumCode)
	if err != nil {
		num = 0
	}

	c.NumCode = num

	cleanValue := strings.ReplaceAll(raw.Value, ",", ".")
	val, err := strconv.ParseFloat(cleanValue, 64)
	if err != nil {
		val = 0
	}

	c.Value = val
	c.CharCode = raw.CharCode
	return nil
}
