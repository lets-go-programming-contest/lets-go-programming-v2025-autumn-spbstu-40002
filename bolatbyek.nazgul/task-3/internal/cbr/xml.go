package cbr

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	XMLName  xml.Name `xml:"Valute"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Value    string   `xml:"Value"`
}

func ParseXML(filePath string) (*ValCurs, error) {
	xmlData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file %q: %w", filePath, err)
	}

	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("failed to decode XML from file %q: %w", filePath, err)
	}

	return &valCurs, nil
}
