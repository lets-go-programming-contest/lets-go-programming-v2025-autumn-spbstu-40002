package xmlread

import (
	"encoding/xml"
	"os"

	"github.com/manyanin.alexander/task-3/internal/currency"
	"github.com/manyanin.alexander/task-3/internal/errors"
	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	XMLName xml.Name            `xml:"ValCurs"`
	Valutes []currency.Currency `xml:"Valute"`
}

func ReadCurrenciesFromXML(filePath string) []currency.Currency {
	file, err := os.Open(filePath)
	if err != nil {
		panic(errors.ErrXMLRead.Error() + ": " + err.Error())
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

	return valCurs.Valutes
}
