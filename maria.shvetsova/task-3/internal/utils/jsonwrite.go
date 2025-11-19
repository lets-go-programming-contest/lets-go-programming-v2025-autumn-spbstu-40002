package utils

import (
	"encoding/json"
	"errors"
	"os"
)

var errEncoding = errors.New("failed to encode")

func WriteToJSON(data []Valute, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return errFileCreating
	}

	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(data)
	if err != nil {
		return errEncoding
	}

	return nil
}
