package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func SaveToJSON(data any, path string) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	bytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	err = os.WriteFile(path, bytes, 0600)
	if err != nil {
		return fmt.Errorf("write json: %w", err)
	}

	return nil
}
