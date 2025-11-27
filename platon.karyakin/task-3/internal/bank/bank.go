package bank

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/charmap"
)

type Bank struct {
	Date       string     `xml:"Date,attr" json:"date"`
	Name       string     `xml:"name,attr" json:"name"`
	Currencies []Currency `xml:"Valute" json:"currencies"`
}

type Currency struct {
	ID       string `xml:"ID,attr" json:"id"`
	NumCode  int    `xml:"NumCode" json:"num_code"`
	CharCode string `xml:"CharCode" json:"char_code"`
	Value    string `xml:"Value" json:"value"`
	Nominal  int    `xml:"Nominal" json:"nominal"`
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
