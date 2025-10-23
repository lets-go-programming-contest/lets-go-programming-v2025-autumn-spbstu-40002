package parsers

import (
	"encoding/xml"
)

// Data type for floating point numbers with comma.
type CommaFloat64 float64

// Data struct for XML/JSON marshalling/unmarshalling.
type ValStruct struct {
	XMLName xml.Name `xml:"ValCurs" json:"-"`
	Text    string   `xml:",chardata" json:"-"`
	Date    string   `xml:"Date,attr" json:"-"`
	Name    string   `xml:"name,attr" json:"-"`
	Valute  []struct {
		Text      string       `xml:",chardata" json:"-"`
		ID        string       `xml:"ID,attr" json:"-"`
		NumCode   int          `xml:"NumCode" json:"num_code"`
		CharCode  string       `xml:"CharCode" json:"char_code"`
		Nominal   string       `xml:"Nominal" json:"-"`
		Name      string       `xml:"Name" json:"-"`
		Value     CommaFloat64 `xml:"Value" json:"value"`
		VunitRate CommaFloat64 `xml:"VunitRate" json:"-"`
	} `xml:"Valute" json:"-"`
}
