package parser

import (
	"encoding/xml"
	"os"

	"github.com/manyanin.alexander/task-3/internal/errors"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

func ParseXML(filePath string) *ValCurs {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		panic(errors.ErrInputFileNotExist.Error() + ": " + filePath)
	}

	xmlFile, err := os.Open(filePath)
	if err != nil {
		panic(errors.ErrXMLRead.Error() + ": " + filePath)
	}
	defer xmlFile.Close()

	var valCurs ValCurs
	decoder := xml.NewDecoder(xmlFile)
	err = decoder.Decode(&valCurs)
	if err != nil {
		panic(errors.ErrXMLDecode.Error() + ": " + err.Error())
	}

	if len(valCurs.Valutes) == 0 {
		panic(errors.ErrXMLEmpty)
	}

	return &valCurs
}
