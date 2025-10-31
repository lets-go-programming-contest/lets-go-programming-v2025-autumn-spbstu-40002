package xml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"sort"
	"golang.org/x/net/html/charset"
)

var ErrNotFound = errors.New("no such file or directory")

type Currencies struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date string `xml:"Date,attr"`
	Name string `xml:"name,attr"`
	Currencies []Currency `xml:"Valute"`
}

func LoadCurrencies(file string) (*Currencies, error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, ErrNotFound
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("close error: %v\n", err)
		}
	}()
	var data Currencies
	dec := xml.NewDecoder(f)
	dec.CharsetReader = charset.NewReaderLabel
	if err := dec.Decode(&data); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &data, nil
}

func (c *Currencies) SortByValue() {
	sort.Sort(ByExchangeRate(c.Currencies))
}
