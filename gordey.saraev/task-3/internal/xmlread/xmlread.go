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
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic(errors.ErrXMLFileRead.Error() + ": " + closeErr.Error())
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	err = decoder.Decode(&valCurs)
	if err != nil {
		panic(errors.ErrXMLDecode.Error() + ": " + err.Error())
	}

	currencies := make([]currency.Currency, 0, len(valCurs.Valutes))

	for _, valute := range valCurs.Valutes {
		currency := currency.ValuteToCurrency(valute)
		currencies = append(currencies, *currency)
	}

	if len(currencies) == 0 {
		panic(errors.ErrNoCurrenciesExtracted.Error())
	}

	return currencies
}
