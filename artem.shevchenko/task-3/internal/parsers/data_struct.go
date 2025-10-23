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
		Text      string       `json:"-"         xml:",chardata"`
		ID        string       `json:"-"         xml:"ID,attr"`
		NumCode   int          `json:"num_code"  xml:"NumCode"`
		CharCode  string       `json:"char_code" xml:"CharCode"`
		Nominal   string       `json:"-"         xml:"Nominal"`
		Name      string       `json:"-"         xml:"Name"`
		Value     CommaFloat64 `json:"value"     xml:"Value"`
		VunitRate CommaFloat64 `json:"-"         xml:"VunitRate"`
	} `json:"-" xml:"Valute"`
}
