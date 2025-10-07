package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func SaveJSON(path string, items []OutputItem) {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			panic(err)
		}
	}
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		panic(err)
	}
	if _, err := f.Write(data); err != nil {
		panic(err)
	}
}
