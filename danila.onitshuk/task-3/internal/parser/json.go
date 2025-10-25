package parser

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const dirPerm = 0o755

type JSONData struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func WriteJSON(path string, recordData []JSONData) {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		panic(ErrCreatDir)
	}

	file, err := os.Create(path)
	if err != nil {
		panic(ErrCreatJSON)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			panic(ErrCloseFile)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(recordData)
	if err != nil {
		panic(ErrWriteJSON)
	}
}
