package parsers

import (
	"encoding/json"
	"os"
	"path/filepath"

	merr "github.com/slendycs/go-lab-3/internal/myerrors"
)

func (value CommaFloat64) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(value))
}

func WriteJSON(path string, data *ValStruct) {
	// Serialize data
	rawData, err := json.MarshalIndent(data.Valute, "", "  ")
	if err != nil {
		panic(merr.ErrFailedToSerializeJSON)
	}

	// Creating output directory
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(merr.ErrFailedToCreateDir)
	}

	// Try to open output file.
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(merr.ErrFailedToOpenOutputFile)
	}
	defer file.Close()

	// Write data
	_, err = file.Write(rawData)
	if err != nil {
		panic(merr.ErrFailedToWriteData)
	}
}
