package writer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/eg0sha-0/task-3/internal/reader"
)

func Save(valutes []reader.Valute, path string) error {
	valutesJSON, err := json.MarshalIndent(valutes, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	const (
		filePermissionCode      = 0o644
		directoryPermissionCode = 0o755
	)

	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, directoryPermissionCode); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	err = os.WriteFile(path, valutesJSON, filePermissionCode)
	if err != nil {
		return fmt.Errorf("failed to write json: %w", err)
	}

	return nil
}
