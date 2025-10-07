package utils

import (
	"encoding/xml"
	"os"

	"golang.org/x/net/html/charset"
)

type cbrValCurs struct {
	XMLName xml.Name   `xml:"ValCurs"`
	Valute  []cbrValue `xml:"Valute"`
}

type cbrValue struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  string `xml:"Nominal"`
	Value    string `xml:"Value"`
}

func ParseCBRXML(path string) []cbrValue {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer func() {
		cerr := file.Close()
		if cerr != nil {
			panic(cerr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var root cbrValCurs
	err = decoder.Decode(&root)
	if err != nil {
		panic(err)
	}

	return root.Valute
}
