package xmlparser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"

	mdls "github.com/ControlShiftEscape/task-3/internal/models"
	"golang.org/x/net/html/charset"
)

var (
	ErrXMLFileOpen  = errors.New("failed to open XML file")
	ErrXMLParsing   = errors.New("failed to parse XML")
	ErrNoCurrencies = errors.New("no currencies found in XML")
)

func ParseXML(path string) (*mdls.ValCurs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: Input file %s: %w", ErrXMLFileOpen, path, err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close XML file: %v\n", closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var curs mdls.ValCurs
	if err := decoder.Decode(&curs); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrXMLParsing, err)
	}

	if len(curs.Valutes) == 0 {
		return nil, ErrNoCurrencies
	}

	return &curs, nil
}
