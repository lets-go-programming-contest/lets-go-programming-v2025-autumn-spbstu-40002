package xmlutils

import (
	"encoding/xml"
	"os"

	"golang.org/x/net/html/charset"
)

type CurrencyList struct {
	Items []Currency `xml:"Valute"`
}

func ReadXML(path string) []Currency {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	var list CurrencyList

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&list); err != nil {
		panic(err)
	}

	return list.Items
}
