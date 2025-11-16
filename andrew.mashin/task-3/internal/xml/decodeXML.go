package xml

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"

	"golang.org/x/net/html/charset"
)

type Currencies struct {
	Currencies []Currency `xml:"Valute"`
}

func GetCurrencies(fileName string) (*Currencies, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	var currencies Currencies

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&currencies)

	return &currencies, nil
}

func (currencies *Currencies) SortOfCurrencies() {
	sort.Sort(ByValue(currencies.Currencies))
}
