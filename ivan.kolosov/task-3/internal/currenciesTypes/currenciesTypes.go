package currenciestypes

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	errParsingXML   = errors.New("error occurred while parsing xml file")
	errParsingFloat = errors.New("error occurred while parsing float")
)

type Currency struct {
	NumCode  int     `xml:"NumCode" json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    float64 `xml:"Value" json:"value"`
}

type Currencies struct {
	CurrencyArray []Currency `xml:"Valute"`
}

func (currency *Currency) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var temp struct {
		NumCode  int    `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	err := dec.DecodeElement(&temp, &start)
	if err != nil {
		return fmt.Errorf("%w: %w", errParsingXML, err)
	}

	str := strings.ReplaceAll(temp.Value, ",", ".")
	if str == "" {
		currency.Value = 0.0
	} else {
		currency.Value, err = strconv.ParseFloat(str, 64)
		if err != nil {
			return fmt.Errorf("%w: %w", errParsingFloat, err)
		}
	}

	currency.NumCode = temp.NumCode
	currency.CharCode = temp.CharCode

	return nil
}
