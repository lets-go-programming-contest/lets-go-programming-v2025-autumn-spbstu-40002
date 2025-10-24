package xmlfiles

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"

	"golang.org/x/net/html/charset"
)

type CurrenciesXML struct {
	XMLName    xml.Name      `xml:"ValCurs"`
	Date       string        `xml:"Date,attr"`
	Name       string        `xml:"name,attr"`
	Currencies []CurrencyXML `xml:"Valute"`
}

func GetCurrencies(fileName string) (*CurrenciesXML, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	var currencies CurrenciesXML

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(&currencies)

	return &currencies, nil
}

func (currencies *CurrenciesXML) SortOfCurrencies() {
	sort.Sort(Currencies(currencies.Currencies))
}
