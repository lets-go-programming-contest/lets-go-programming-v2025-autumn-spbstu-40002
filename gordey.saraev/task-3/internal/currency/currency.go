package currency

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type Currency struct {
	XMLName  xml.Name `xml:"Valute" json:"-"`
	NumCode  int      `json:"num_code"`
	CharCode string   `json:"char_code"`
	Value    float64  `json:"value"`
}

func (c *Currency) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var raw struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	if err := d.DecodeElement(&raw, &start); err != nil {

		return err
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
