package utils

import (
	"encoding/xml"
	"fmt"
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
		return nil, fmt.Errorf("open xml file: %w", err)
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var exch Exchange
	err = decoder.Decode(&exch)
	if err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	return &exch, nil
}
