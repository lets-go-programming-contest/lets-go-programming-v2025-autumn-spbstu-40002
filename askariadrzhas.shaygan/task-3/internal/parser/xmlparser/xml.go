package xmlparser

import (
	"encoding/xml"
	"os"

	"github.com/XShaygaND/task-3/internal/currency"
	"github.com/XShaygaND/task-3/internal/myerrors"
	"github.com/XShaygaND/task-3/internal/utils"
	"golang.org/x/net/html/charset"
)

func ParseXML(path string) *currency.ValCurs {
	file, err := os.Open(path)
	if err != nil {
		panic(myerrors.ErrFileNotFound)
	}
	defer file.Close()

	var valCurs currency.ValCurs
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&valCurs); err != nil {
		panic(myerrors.ErrXMLDecode)
	}

	utils.SortValutesByValue(valCurs.Valute)
	return &valCurs
}
