package xml

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"

	"golang.org/x/net/html/charset"
)

type Currencies struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Date       string     `xml:"Date,attr"`
	Name       string     `xml:"name,attr"`
	Currencies []Currency `xml:"Valute"`
}

func LoadCurrencies(filename string) (*Currencies, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("no such file or directory")
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic(closeErr)
		}
	}()

	var currencies Currencies
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&currencies); err != nil {
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}

	return &currencies, nil
}

func (c *Currencies) SortByValue() {
	sort.Sort(ByExchangeRate(c.Currencies))
}
