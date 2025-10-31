package xmlread

import (
	"encoding/xml"
	"os"

	currency "github.com/F0LY/task-3/internal/currency"
	errors "github.com/F0LY/task-3/internal/errors"
	"golang.org/x/net/html/charset"
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

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic(errors.ErrXMLFileRead.Error() + ": " + closeErr.Error())
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	if err = decoder.Decode(&valCurs); err != nil {
		panic(errors.ErrXMLDecode.Error() + ": " + err.Error())
	}

	currencies := make([]currency.Currency, 0, len(valCurs.Valutes))

	for _, valute := range valCurs.Valutes {
		curr := currency.ValuteToCurrency(valute)
		currencies = append(currencies, *curr)
	}

	if len(currencies) == 0 {
		panic(errors.ErrNoCurrenciesExtracted.Error())
	}

	return currencies
}
