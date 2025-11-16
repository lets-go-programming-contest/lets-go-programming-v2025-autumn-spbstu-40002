package writingcurrencies

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	cur "github.com/InsomniaDemon/task-3/internal/currenciesTypes"
)

var (
	errCreatingDir     = errors.New("error occurred while creating specified directory")
	errCreatingFile    = errors.New("error occurred while creating specified file")
	errClosingJSONFile = errors.New("error occurred while closing json file")
	errEncodingJSON    = errors.New("error occurred while encoding json file")
)

const dirPerm = 0o755

func WriteCurrencies(data cur.Currencies, path string) (returnError error) {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		return fmt.Errorf("%w: %w", errCreatingDir, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%w: %w", errCreatingFile, err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			if returnError != nil {
				returnError = fmt.Errorf("%w: %w; %w", returnError, err, errClosingJSONFile)
			} else {
				returnError = fmt.Errorf("%w: %w", errClosingJSONFile, err)
			}
		}
	}()

	enc := json.NewEncoder(file)
	enc.SetIndent("", " ")

	err = enc.Encode(data.CurrencyArray)
	if err != nil {
		return fmt.Errorf("%w: %w", errEncodingJSON, err)
	}

	return nil
}
