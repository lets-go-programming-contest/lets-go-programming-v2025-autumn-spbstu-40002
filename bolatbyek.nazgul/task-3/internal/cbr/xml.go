package cbr

import (
	"encoding/xml"
	"os"
)

// Currency represents a currency entry
type Currency struct {
	NumCode  int     `xml:"NumCode" json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    float64 `xml:"Value" json:"value"`
}

// ValCurs represents the XML structure from CBR
type ValCurs struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Date       string     `xml:"Date,attr"`
	Name       string     `xml:"Name,attr"`
	Currencies []Currency `xml:"Valute"`
}

// Parser handles XML parsing for CBR data
type Parser struct{}

// NewParser creates a new XML parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseXML parses XML data from file
func (p *Parser) ParseXML(filename string) (*ValCurs, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic("Failed to read input file: " + err.Error())
	}

	var valCurs ValCurs
	err = xml.Unmarshal(data, &valCurs)
	if err != nil {
		panic("Failed to parse XML: " + err.Error())
	}

	return &valCurs, nil
}