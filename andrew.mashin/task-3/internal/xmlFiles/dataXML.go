package xmlFiles

import (
	"fmt"
	"strconv"
	"strings"
)

type CurrencyXML struct {
	ID        string `xml:"ID,attr"`
	NumCode   int    `xml:"NumCode"`
	CharCode  string `xml:"CharCode"`
	Nominal   int    `xml:"Nominal"`
	Name      string `xml:"Name"`
	Value     string `xml:"Value"`
	VunitRate string `xml:"VunitRate"`
}

func (currency CurrencyXML) GetFloat() (float64, error) {
	commaReplacement := strings.Replace(currency.Value, ",", ".", -1)

	value, err := strconv.ParseFloat(commaReplacement, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert value %s to float", commaReplacement)
	}

	return value, nil
}

type Currencies []CurrencyXML

func (currency Currencies) Len() int {
	return len(currency)
}

func (currency Currencies) Swap(i, j int) {
	currency[i], currency[j] = currency[j], currency[i]
}

func (currency Currencies) Less(i, j int) bool {
	currencyI, err := currency[i].GetFloat()
	if err != nil {
		panic(err)
	}

	currencyJ, err := currency[j].GetFloat()
	if err != nil {
		panic(err)
	}

	return currencyI > currencyJ
}
