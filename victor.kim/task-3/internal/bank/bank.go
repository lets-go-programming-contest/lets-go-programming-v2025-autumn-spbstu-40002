package bank

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/victor.kim/task-3/pkg/must"
	"golang.org/x/text/encoding/charmap"
)

type Currency struct {
	ID       string `xml:"ID,attr" json:"-"`
	NumCode  int    `xml:"NumCode"    json:"num_code"`
	CharCode string `xml:"CharCode"   json:"char_code"`
	Value    string `xml:"Value"      json:"value_raw"`
}

type Bank struct {
	Date       string     `xml:"Date,attr" json:"date"`
	Name       string     `xml:"name,attr" json:"name"`
	Currencies []Currency `xml:"Valute"    json:"currencies"`
}

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch charset {
	case "windows-1251":
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	}

	return input, nil
}

func ParseXML(r io.Reader) (*Bank, error) {
	dec := xml.NewDecoder(r)
	dec.CharsetReader = charsetReader

	out := new(Bank)
	if err := dec.Decode(out); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	return out, nil
}

func ParseFileXML(path string) (*Bank, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open input: %w", err)
	}

	defer must.Close(path, f)

	return ParseXML(f)
}
