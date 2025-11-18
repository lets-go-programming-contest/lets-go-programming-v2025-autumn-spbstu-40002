package jsonwriter

import (
	"encoding/json"
	"fmt"
	"os"
)

func SaveToJSON(currencies []interface{}, filePath string) error {
	outputJSON, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal currencies to JSON: %w", err)
	}

	err = os.WriteFile(filePath, outputJSON, 0600)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file %q: %w", filePath, err)
	}

	return nil
}
