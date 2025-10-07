package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func SaveToJSON(data any, path string) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0644)
}
