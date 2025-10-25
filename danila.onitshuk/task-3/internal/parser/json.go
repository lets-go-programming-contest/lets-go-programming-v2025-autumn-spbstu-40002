package parser

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type JsonData struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func WriteJson(path string, recordData []JsonData) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(ErrCreatDir)
	}

	file, err := os.Create(path)
	if err != nil {
		panic(ErrCreatJson)
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
		panic(ErrWriteJson)
	}
}
