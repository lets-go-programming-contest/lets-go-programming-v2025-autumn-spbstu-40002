package parser

import (
	"encoding/xml"
	"os"

	"github.com/XShaygaND/task-3/internal/parser/myerrors"
)

type Item struct {
	ID    string `xml:"id"`
	Value string `xml:"value"`
}

func ParseXML(filePath string) ([]Item, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var items []Item
	err = xml.Unmarshal(data, &items)
	if err != nil {
		return nil, myerrors.ErrParseXML
	}

	return items, nil
}
