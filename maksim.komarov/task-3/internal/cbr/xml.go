package cbr

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
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
	f, err := os.Open(path)
	if err != nil {
		return Document{}, fmt.Errorf("%w: %v", ErrOpenInputXML, err)
	}
	defer f.Close()

	var d Document
	dec := xml.NewDecoder(f)
	if err := dec.Decode(&d); err != nil {
		return Document{}, fmt.Errorf("%w: %v", ErrDecodeInputXML, err)
	}

	return d, nil
}
