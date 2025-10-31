package xml

import (
	"encoding/xml"
	"os"

	"github.com/rachguta/task-3/internal/currency"
	"github.com/rachguta/task-3/internal/myerrors"
	"github.com/rachguta/task-3/internal/utils"
	"golang.org/x/net/html/charset"
)

func ParseXML(path string) *currency.ValCurs {
	file, err := os.Open(path)
	if err != nil {
		panic(myerrors.ErrFileNotFound)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			panic(myerrors.ErrCloseFile)
		}
	}()

	var valCurs currency.ValCurs

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(&valCurs)
	if err != nil {
		panic(myerrors.ErrXMLDecode)
	}

	utils.SortValutesByValue(valCurs.Valute)

	return &valCurs
}
