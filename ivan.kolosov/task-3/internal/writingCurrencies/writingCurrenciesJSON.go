package writingсurrencies

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	errCreatingDir     = errors.New("error occurred while creating specified directory")
	errCreatingFile    = errors.New("error occurred while creating specified file")
	errClosingJSONFile = errors.New("error occurred while closing json file")
	errEncodingJSON    = errors.New("error occurred while encoding json file")
)

const dirPerm = 0o755

func WriteCurrencies(data CurrenciesJSON, path string) (returnError error) {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		return fmt.Errorf("%v: %w", errCreatingDir, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%v: %w", errCreatingFile, err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			if returnError != nil {
				returnError = fmt.Errorf("%v: %w; %v", returnError, err, errClosingJSONFile)
			} else {
				returnError = fmt.Errorf("%v: %w", errClosingJSONFile, err)
			}
		}
	}()

	enc := json.NewEncoder(file)
	enc.SetIndent("", " ")

	err = enc.Encode(data)
	if err != nil {
		return fmt.Errorf("%v: %w", errEncodingJSON, err)
	}

	return nil
}
