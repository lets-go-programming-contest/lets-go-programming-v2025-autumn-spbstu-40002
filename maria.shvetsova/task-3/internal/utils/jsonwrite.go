package utils

import (
	"encoding/json"
	"errors"
	"os"
)

var (
	errInvalidFormat = errors.New("invalid format")
	errEncoding      = errors.New("failed to encode")
)

type Output struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func WriteToJSON(data []Output, outputPath string) error {
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
