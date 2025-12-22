package currency

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type ExchangeRate float64

type CurrencyIndex struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	NumCode  int          `json:"num_code"  xml:"NumCode"`
	CharCode string       `json:"char_code" xml:"CharCode"`
	Value    ExchangeRate `json:"value"     xml:"Value"`
}

func (exchangeRate *ExchangeRate) UnmarshalXML(decoder *xml.Decoder, startElement xml.StartElement) error {
	var str string

	err := decoder.DecodeElement(&str, &startElement)
	if err != nil {
		return fmt.Errorf("decode element: %w", err)
	}

	str = strings.Replace(str, ",", ".", 1)

	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("parse float: %w", err)
	}

	*exchangeRate = ExchangeRate(value)

	return nil
}
