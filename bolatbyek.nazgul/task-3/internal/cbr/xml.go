package cbr

import (
	"encoding/xml"
	"io/ioutil"
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
	xmlData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var valCurs ValCurs
	err = xml.Unmarshal(xmlData, &valCurs)
	if err != nil {
		return nil, err
	}

	return &valCurs, nil
}
