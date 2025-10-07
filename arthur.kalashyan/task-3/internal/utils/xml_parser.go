package utils

import (
	"encoding/xml"
	"io"
	"os"
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
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var root cbrValCurs
	if err := xml.Unmarshal(data, &root); err != nil {
		panic(err)
	}
	return root.Valute
}
