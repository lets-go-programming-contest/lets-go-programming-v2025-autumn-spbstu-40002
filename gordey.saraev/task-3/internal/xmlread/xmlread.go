package xmlread

import (
	"encoding/xml"
	"os"

	"golang.org/x/net/html/charset"

	"github.com/F0LY/task-3/internal/currency"
	"github.com/F0LY/task-3/internal/errors"
)

type ValCurs struct {
	XMLName xml.Name          `xml:"ValCurs"`
	Valutes []currency.Valute `xml:"Valute"`
}

func ReadCurrenciesFromXML(filePath string) []currency.Currency {
	file, err := os.Open(filePath)
	if err != nil {
		panic(errors.ErrXMLFileRead.Error() + ": " + err.Error())
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	err = decoder.Decode(&valCurs)
	if err != nil {
		panic(errors.ErrXMLDecode.Error() + ": " + err.Error())
	}

	var currencies []currency.Currency
	for _, valute := range valCurs.Valutes {
		currency := currency.ValuteToCurrency(valute)
		if currency != nil {
			currencies = append(currencies, *currency)
		}
	}

	return currencies
}
