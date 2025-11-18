package xmlparser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/Nazkaaa/task-3/internal/models"
	"golang.org/x/net/html/charset"
)

func ParseXML(filePath string) (*models.ValCurs, error) {
	xmlData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file %q: %w", filePath, err)
	}

	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs models.ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("failed to decode XML from file %q: %w", filePath, err)
	}

	return &valCurs, nil
}
