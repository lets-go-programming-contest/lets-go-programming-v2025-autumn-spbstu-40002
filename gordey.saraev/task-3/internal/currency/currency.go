package currency

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type Currency struct {
	XMLName  xml.Name `json:"-"         xml:"Valute"`
	NumCode  int      `json:"num_code"  xml:"NumCode"`
	CharCode string   `json:"char_code" xml:"CharCode"`
	Value    float64  `json:"value"     xml:"Value"`
}

func (c *Currency) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var raw struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	err := decoder.DecodeElement(&raw, &start)
	if err != nil {
		return fmt.Errorf("decode element: %w", err)
	}

	num, err := strconv.Atoi(raw.NumCode)
	if err != nil {
		c.NumCode = 0
	} else {
		c.NumCode = num
	}

	cleanValue := strings.ReplaceAll(raw.Value, ",", ".")
	val, err := strconv.ParseFloat(cleanValue, 64)

	if err != nil {
		c.Value = 0
	} else {
		c.Value = val
	}

	c.CharCode = raw.CharCode

	return nil
}
