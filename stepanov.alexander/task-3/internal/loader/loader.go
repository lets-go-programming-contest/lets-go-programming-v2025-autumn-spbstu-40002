package loader

import (
	"encoding/xml"
	"os"

	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

func LoadXML(path string) (*ValCurs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs

	if err := decoder.Decode(&valCurs); err != nil {
		return nil, err
	}

	return &valCurs, nil
}
