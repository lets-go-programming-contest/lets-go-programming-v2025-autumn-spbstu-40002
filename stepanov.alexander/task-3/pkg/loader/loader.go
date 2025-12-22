package loader

import (
	"encoding/xml"
	"os"
)

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode   string `xml:"NumCode"`
	CharCode  string `xml:"CharCode"`
	Value     string `xml:"Value"`
}

func LoadXML(filepath string) (*ValCurs, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var valCurs ValCurs

	err = xml.Unmarshal(data, &valCurs)
	if err != nil {
		return nil, err
	}

	return &valCurs, nil
}
