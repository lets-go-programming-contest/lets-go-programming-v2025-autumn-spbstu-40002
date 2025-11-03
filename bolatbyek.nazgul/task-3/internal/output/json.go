package output

import (
	"encoding/json"
	"os"
)

func SaveToJSON(currencies []interface{}, filePath string) error {
	outputJSON, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, outputJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}
