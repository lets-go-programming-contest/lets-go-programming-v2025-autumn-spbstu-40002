package bank

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/Danya-byte/task-3/pkg/must"
	"golang.org/x/text/encoding/charmap"
)

type Bank struct {
	Date       string     `json:"date"       xml:"Date,attr"`
	Name       string     `json:"name"       xml:"name,attr"`
	Currencies []Currency `json:"currencies" xml:"Valute"`
}

type Currency struct {
	ID       string `json:"id"        xml:"ID,attr"`
	NumCode  int    `json:"num_code"  xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    string `json:"value"     xml:"Value"`
	Nominal  int    `json:"nominal"   xml:"Nominal"`
	Name     string `json:"name"      xml:"Name"`
	Rate     string `json:"rate"      xml:"VunitRate"`
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

	bank := new(Bank)

	if err := decoder.Decode(&bank); err != nil {
		return nil, fmt.Errorf("decoding currency bank: %w", err)
	}

	return bank, nil
}

func ParseFileXML(path string) (*Bank, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open input file: %w", err)
	}

	defer must.Close(path, file)

	return ParseXML(file)
}
