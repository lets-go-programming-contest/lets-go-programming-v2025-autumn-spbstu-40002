package xmlparse

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"

	model "github.com/HuaChenju/task-3/internal/xmlfile/model"
	"golang.org/x/net/html/charset"
)

var (
	errOpeningFile = errors.New("error with open XML file")
	errXMLParsing  = errors.New("error with XML parsing")
	errClosingFile = errors.New("error with closing XML file")
)

func GetValCursStruct(inputPath string) (_ model.ValCurs, returnError error) {
	var doc model.ValCurs

	file, err := os.Open(inputPath)
	if err != nil {
		return doc, fmt.Errorf("%w: %w", errOpeningFile, err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			if returnError != nil {
				returnError = fmt.Errorf("%w: %w; %w", returnError, err, errClosingFile)
			} else {
				returnError = fmt.Errorf("%w: %w", errClosingFile, err)
			}
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&doc); err != nil {
		return doc, fmt.Errorf("%w: %w", errXMLParsing, err)
	}

	return doc, nil
}
