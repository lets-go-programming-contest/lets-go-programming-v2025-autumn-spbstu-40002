package reader

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type Valute struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value"     xml:"Value"`
}

type Exchange struct {
	Valutes []Valute `xml:"Valute"`
}

func (v *Valute) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	type Alias Valute

	temporaryStruct := &struct {
		Value string `xml:"Value"`
		*Alias
	}{
		Value: "",
		Alias: (*Alias)(v),
	}

	if err := decoder.DecodeElement(temporaryStruct, &start); err != nil {
		return fmt.Errorf("failed to decode valute: %w", err)
	}

	valueStr := strings.ReplaceAll(temporaryStruct.Value, ",", ".")

	floatValue, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fmt.Errorf("failed to parse float for %s: %w", temporaryStruct.Value, err)
	}

	v.Value = floatValue

	return nil
}
