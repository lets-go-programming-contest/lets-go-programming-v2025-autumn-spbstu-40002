package bank

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Danya-byte/task-3/pkg/must"
	"golang.org/x/text/encoding/charmap"
)

type Bank struct {
	Date       string     `xml:"Date,attr" json:"date"`
	Name       string     `xml:"name,attr" json:"name"`
	Currencies []Currency `xml:"Valute" json:"currencies"`
}

type Currency struct {
	ID       string  `xml:"ID,attr" json:"id"`
	NumCode  int     `xml:"NumCode" json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    float64 `json:"value"`
	Nominal  int     `xml:"Nominal" json:"nominal"`
	Name     string  `xml:"Name" json:"name"`
	Rate     float64 `json:"rate"`
}

func (c *Currency) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type raw Currency // избегаем рекурсии
	var tmp struct {
		raw
		Value string `xml:"Value"`
		Rate  string `xml:"VunitRate"`
	}

	if err := d.DecodeElement(&tmp, &start); err != nil {
		return err
	}

	*c = Currency(tmp.raw)

	if v := strings.Replace(tmp.Value, ",", ".", 1); v != "" {
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("invalid Value: %w", err)
		}
		c.Value = val
	}

	if v := strings.Replace(tmp.Rate, ",", ".", 1); v != "" {
		rval, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("invalid Rate: %w", err)
		}
		c.Rate = rval
	}

	return nil
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
	if err := decoder.Decode(bank); err != nil {
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
