package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func writeJSON(out []OutputCurrency, path string) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return fmt.Errorf("%v: %w", ErrDirCreate, err)
		}
	}

	bytes, err := json.MarshalIndent(out, "", " ")
	if err != nil {
		return fmt.Errorf("%v: %w", ErrJSONWrite, err)
	}

	if err = os.WriteFile(path, bytes, filePerm); err != nil {
		return fmt.Errorf("%v: %w", ErrJSONWrite, err)
	}

	return nil
}
