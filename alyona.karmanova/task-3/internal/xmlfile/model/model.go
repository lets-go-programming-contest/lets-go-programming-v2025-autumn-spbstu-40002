package model

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type ValCurs struct {
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int     `json:"num_code"  xml:"num_code"  yaml:"num_code"`
	CharCode string  `json:"char_code" xml:"char_code" yaml:"char_code"`
	Value    float64 `json:"value"     xml:"value"     yaml:"value"`
}

func (v *Valute) UnmarshalXML(decod *xml.Decoder, start xml.StartElement) error {
	var aux struct {
		NumCode  int    `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	if err := decod.DecodeElement(&aux, &start); err != nil {
		return fmt.Errorf("couldn't decode elem: %w", err)
	}

	valStr := strings.ReplaceAll(strings.TrimSpace(aux.Value), ",", ".")

	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return fmt.Errorf("couldn't parse Value: %w", err)
	}

	v.NumCode = aux.NumCode
	v.CharCode = aux.CharCode
	v.Value = val

	return nil
}
