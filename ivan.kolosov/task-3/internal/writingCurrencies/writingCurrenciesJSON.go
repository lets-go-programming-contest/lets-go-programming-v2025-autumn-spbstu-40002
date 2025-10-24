package writingcurrencies

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	errCreatingDir     = "error occurred while creating specified directory"
	errCreatingFile    = "error occurred while creating specified file"
	errClosingJSONFile = "error occurred while closing json file"
	errEncodingJSON    = "error occurred while encoding json file"
)

const dirPerm = 0o755

func WriteCurrencies(data CurrenciesJSON, path string) (returnError error) {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		return fmt.Errorf("%s: %w", errCreatingDir, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%s: %w", errCreatingFile, err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			if returnError != nil {
				returnError = fmt.Errorf("%s: %w; %s", returnError, err, errClosingJSONFile)
			} else {
				returnError = fmt.Errorf("%s: %w", errClosingJSONFile, err)
			}
		}
	}()

	enc := json.NewEncoder(file)
	enc.SetIndent("", " ")

	err = enc.Encode(data)
	if err != nil {
		return fmt.Errorf("%s: %w", errEncodingJSON, err)
	}

	return nil
}
