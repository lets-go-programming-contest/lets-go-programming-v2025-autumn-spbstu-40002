package parsers

import (
	"encoding/json"
	"os"
	"path/filepath"

	merr "github.com/slendycs/go-lab-3/internal/myerrors"
)

const (
	defaultDirPerm  = 0o755
	defaultFilePerm = 0o644
)

func (value CommaFloat64) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(float64(value))
	if err != nil {
		panic(err)
	}

	return data, nil
}

func WriteJSON(path string, data *ValStruct) {
	// Serialize data.
	rawData, err := json.MarshalIndent(data.Valute, "", "  ")
	if err != nil {
		panic(merr.ErrFailedToSerializeJSON)
	}

	// Creating output directory.
	dir := filepath.Dir(path)

	err = os.MkdirAll(dir, defaultDirPerm)
	if err != nil {
		panic(merr.ErrFailedToCreateDir)
	}

	// Try to open output file.
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, defaultFilePerm)
	if err != nil {
		panic(merr.ErrFailedToOpenOutputFile)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			panic(merr.ErrFailedToCloseFile)
		}
	}()

	// Write data.
	_, err = file.Write(rawData)
	if err != nil {
		panic(merr.ErrFailedToWriteData)
	}
}
