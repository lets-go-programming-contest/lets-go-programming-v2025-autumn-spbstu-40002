package xmlparse

import (
	"encoding/xml"
	"fmt"
	"os"

	model "github.com/HuaChenju/task-3/internal/xmlfile/model"
	"golang.org/x/net/html/charset"
)

func GetValCursStruct(inputPath string) (model.ValCurs, error) {
	var doc model.ValCurs

	file, err := os.Open(inputPath)
	if err != nil {
		return doc, fmt.Errorf("couldn't open XML file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("error closing file: %v", err)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&doc); err != nil {
		return doc, fmt.Errorf("xml parsing error: %w", err)
	}

	return doc, nil
}
