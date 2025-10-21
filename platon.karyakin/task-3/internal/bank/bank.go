package bank

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/charmap"
)

type Bank struct {
	Date       string     `xml:"Date,attr"`
	Name       string     `xml:"name,attr"`
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	ID       string `xml:"ID,attr"`
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
}

func charsetReader(charset string, reader io.Reader) (io.Reader, error) {
	switch charset {
	case "windows-1251":
		return charmap.Windows1251.NewDecoder().Reader(reader), nil
	default:
		return reader, nil
	}
}

func DecodeXML(reader io.Reader) (*Bank, error) {
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charsetReader

	bankData := new(Bank)
	if err := decoder.Decode(bankData); err != nil {
		return nil, fmt.Errorf("failed to decode XML to Bank struct: %w", err)
	}

	return bankData, nil
}

func LoadFromXML(path string) (*Bank, error) {
	fileHandle, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open input XML file: %w", err)
	}

	defer func() { _ = fileHandle.Close() }()

	return DecodeXML(fileHandle)
}
