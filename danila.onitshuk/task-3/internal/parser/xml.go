package parser

import (
	"encoding/xml"
	"os"

	"golang.org/x/net/html/charset"
)

type Valute struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

type XMLData struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valute  []Valute `xml:"Valute"`
}

func ReadXML(path string, data *XMLData) {
	file, err := os.Open(path)
	if err != nil {
		panic(ErrFileNotFound)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			panic(ErrCloseFile)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(data)
	if err != nil {
		panic(ErrDecodeXML)
	}
}
