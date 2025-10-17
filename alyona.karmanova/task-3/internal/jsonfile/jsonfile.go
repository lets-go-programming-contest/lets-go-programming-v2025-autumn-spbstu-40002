package jsonfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	xmlfile "github.com/HuaChenju/task-3/internal/xmlfile"
)

func ensureOutputDir(path string) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("couldn't create a directory %s: %w", dir, err)
	}

	return nil
}

func WriteJSONToFile(filePath string, doc xmlfile.ValCurs) error {
	if err := ensureOutputDir(filePath); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(doc.Valutes, "", "  ")
	if err != nil {
		return fmt.Errorf("couldn't encode in JSON: %w", err)
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("couldn't write to a file: %w", err)
	}

	return nil
}
