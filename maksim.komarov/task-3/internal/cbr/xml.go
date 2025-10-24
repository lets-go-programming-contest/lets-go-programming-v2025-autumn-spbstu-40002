package cbr

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

var (
	ErrOpenInputXML   = errors.New("open input xml")
	ErrDecodeInputXML = errors.New("decode input xml")
)

type Document struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

func ReadFile(path string) (Document, error) {
	file, err := os.Open(path)
	if err != nil {
		return Document{}, fmt.Errorf("%s: %w", ErrOpenInputXML, err)
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var doc Document

	if err := decoder.Decode(&doc); err != nil {
		return Document{}, fmt.Errorf("%s: %w", ErrDecodeInputXML, err)
	}

	return doc, nil
}
