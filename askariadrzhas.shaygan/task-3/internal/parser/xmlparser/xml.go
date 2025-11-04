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

	defer func() {
		if err := file.Close(); err != nil {
			panic(myerrors.ErrCloseFile)
		}
	}()

	var vc currency.ValCurs

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(&vc)
	if err != nil {
		panic(myerrors.ErrXMLDecode)
	}

	utils.SortValutesByValue(vc.Valute)

	return &vc
}
