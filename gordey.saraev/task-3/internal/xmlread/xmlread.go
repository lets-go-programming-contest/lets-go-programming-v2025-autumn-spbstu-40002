package xmlread

import (
	"encoding/xml"
	"os"

	"github.com/F0LY/task-3/internal/currency"
	"github.com/F0LY/task-3/internal/errors"
	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	XMLName xml.Name            `xml:"ValCurs"`
	Valutes []currency.Currency `xml:"Valute"`
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

	if len(valCurs.Valutes) == 0 {
		panic(errors.ErrNoCurrenciesExtracted.Error())
	}

	return valCurs.Valutes
}
