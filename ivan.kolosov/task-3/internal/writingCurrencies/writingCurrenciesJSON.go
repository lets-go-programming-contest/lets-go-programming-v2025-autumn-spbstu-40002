package writingCurrencies

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
		return errCreatingDir
	}

	file, err := os.Create(path)
	if err != nil {
		return errCreatingFile
	}

	defer func() {
		err := file.Close()
		if err != nil {
			if returnError != nil {
				returnError = fmt.Errorf("%w; %v", returnError, errClosingJSONFile)
			} else {
				returnError = errClosingJSONFile
			}
		}
	}()

	enc := json.NewEncoder(file)
	enc.SetIndent("", " ")

	err = enc.Encode(data)
	if err != nil {
		return errEncodingJSON
	}

	return nil
}
