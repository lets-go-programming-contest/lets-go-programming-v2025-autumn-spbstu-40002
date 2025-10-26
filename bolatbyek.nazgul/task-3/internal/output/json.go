package output

import (
	"encoding/json"
	"io/ioutil"
)

func SaveToJSON(currencies []interface{}, filePath string) error {
	outputJSON, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, outputJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}