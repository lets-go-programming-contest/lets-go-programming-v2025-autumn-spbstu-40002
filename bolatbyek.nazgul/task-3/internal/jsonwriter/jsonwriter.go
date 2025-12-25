package jsonwriter

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	filePerm = 0o600
)

func SaveToJSON(currencies []interface{}, filePath string) error {
	outputJSON, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal currencies to JSON: %w", err)
	}

	err = os.WriteFile(filePath, outputJSON, filePerm)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file %q: %w", filePath, err)
	}

	return nil
}
