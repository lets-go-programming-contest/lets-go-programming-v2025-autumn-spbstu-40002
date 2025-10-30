package jsonfile

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"

	outputfile "github.com/HuaChenju/task-3/internal/checkoutput"
	xmlmodel "github.com/HuaChenju/task-3/internal/xmlfile/model"
	"gopkg.in/yaml.v3"
)

const filePerm = 0o600

func WriteToFile(filePath string, doc xmlmodel.ValCurs, format string) error {
	if err := outputfile.EnsureOutputDir(filePath); err != nil {
		return fmt.Errorf("trouble with JSON: %w", err)
	}

	var (
		data []byte
		err  error
	)

	switch format {
	case "yaml":
		data, err = yaml.Marshal(doc.Valutes)
	case "xml":
		data, err = xml.MarshalIndent(doc.Valutes, "", "  ")
	default:
		data, err = json.MarshalIndent(doc.Valutes, "", "  ")
	}

	if err != nil {
		return fmt.Errorf("couldn't encode in %s: %w", format, err)
	}

	if err := os.WriteFile(filePath, data, filePerm); err != nil {
		return fmt.Errorf("couldn't write to a file: %w", err)
	}

	return nil
}
