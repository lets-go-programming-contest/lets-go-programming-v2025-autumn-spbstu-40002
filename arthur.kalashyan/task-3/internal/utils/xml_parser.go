package utils

import (
	"encoding/xml"
	"os"

	"golang.org/x/net/html/charset"
)

type Exchange struct {
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

func ReadXML(path string) (*Exchange, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var e Exchange
	if err := decoder.Decode(&e); err != nil {
		return nil, err
	}
	return &e, nil
}
