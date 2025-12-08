package parsers

import (
	"encoding/xml"
)

// Data type for floating point numbers with comma.
type CommaFloat64 float64

// Data struct for XML/JSON marshalling/unmarshalling.
type ValStruct struct {
	XMLName xml.Name `json:"-" xml:"ValCurs"`
	Text    string   `json:"-" xml:",chardata"`
	Date    string   `json:"-" xml:"Date,attr"`
	Name    string   `json:"-" xml:"name,attr"`
	Valute  []struct {
		NumCode  int          `json:"num_code"  xml:"NumCode"`
		CharCode string       `json:"char_code" xml:"CharCode"`
		Value    CommaFloat64 `json:"value"     xml:"Value"`
	} `json:"-" xml:"Valute"`
}
