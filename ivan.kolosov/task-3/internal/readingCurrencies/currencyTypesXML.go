package readingCurrencies

import "encoding/xml"

type CurrencyXML struct {
	ID             string `xml:"ID,attr"`
	NumCode        int    `xml:"NumCode"`
	CharCode       string `xml:"CharCode"`
	Nominal        int    `xml:"Nominal"`
	Name           string `xml:"Name"`
	Value          string `xml:"Value"`
	VunitRate      string `xml:"VunitRate"`
	ValueFloat     float64
	VunitRateFloat float64
}

type CurrenciesXML struct {
	XMLName    xml.Name      `xml:"ValCurs"`
	Date       string        `xml:"Date,attr"`
	Name       string        `xml:"name,attr"`
	Currencies []CurrencyXML `xml:"Valute"`
}
