package utils

import (
	"encoding/xml"
	"errors"
	"strconv"
	"strings"
)

var (
	errDecodingElement = errors.New("failed to decode XML element")
	errParsingValue    = errors.New("failed to parse value")
)

type Valute struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value"     xml:"Value"`
}

func (v *Valute) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var aux struct {
		NumCode  int    `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	if err := decoder.DecodeElement(&aux, &start); err != nil {
		return errDecodingElement
	}

	normalized := strings.ReplaceAll(aux.Value, ",", ".")

	value, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return errParsingValue
	}

	v.NumCode = aux.NumCode
	v.CharCode = aux.CharCode
	v.Value = value

	return nil
}

type Valutes []Valute

func (v Valutes) Len() int { return len(v) }

func (v Valutes) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

func (v Valutes) Less(i, j int) bool {
	return v[i].Value > v[j].Value
}
