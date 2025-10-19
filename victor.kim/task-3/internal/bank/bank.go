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

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch charset {
	case "windows-1251":
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	default:
		return input, nil
	}
}

func ParseXML(r io.Reader) (*Bank, error) {
	decoder := xml.NewDecoder(r)
	decoder.CharsetReader = charsetReader

	b := new(Bank)
	if err := decoder.Decode(b); err != nil {
		return nil, fmt.Errorf("decoding currency bank: %w", err)
	}

	return b, nil
}

func ParseFileXML(path string) (*Bank, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open input file: %w", err)
	}

	defer func() { _ = f.Close() }()

	return ParseXML(f)
}
