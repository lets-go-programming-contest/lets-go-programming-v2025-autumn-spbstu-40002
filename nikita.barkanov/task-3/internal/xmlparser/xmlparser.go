package xmlparser

import (
	"encoding/xml"
	"fmt"
	"os"

	mdls "github.com/ControlShiftEscape/task-3/internal/models"
)

var (
	ErrXMLFileOpen = fmt.Errorf("failed to open XML file")
	ErrXMLParsing  = fmt.Errorf("failed to parse XML")
)

func ParseXML(path string) (*mdls.ValCurs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrXMLFileOpen, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close XML file: %v\n", closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	var curs mdls.ValCurs
	if err := decoder.Decode(&curs); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrXMLParsing, err)
	}

	if len(curs.Valutes) == 0 {
		return nil, fmt.Errorf("no currencies found in XML")
	}

	return &curs, nil
}
