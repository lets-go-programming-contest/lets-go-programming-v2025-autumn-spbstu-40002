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
	XMLName    xml.Name   `xml:"ValCurs"`
	Date       string     `xml:"Date,attr"`
	Name       string     `xml:"name,attr"`
	Currencies []Currency `xml:"Valute"`
}

func LoadCurrencies(filename string) (*Currencies, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, ErrNotFound
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("warning: failed to close file: %v\n", closeErr)
		}
	}()

	var currencies Currencies

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&currencies); err != nil {
		return nil, fmt.Errorf("decode XML: %w", err)
	}

	return &currencies, nil
}

func (c *Currencies) SortByValue() {
	sort.Sort(ByExchangeRate(c.Currencies))
}
