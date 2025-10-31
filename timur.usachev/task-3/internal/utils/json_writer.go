package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func writeJSON(out []OutputCurrency, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("%w: %v", ErrDirCreate, err)
	}
	bytes, err := json.MarshalIndent(out, "", " ")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrJSONWrite, err)
	}
	if err := os.WriteFile(path, bytes, filePerm); err != nil {
		return fmt.Errorf("%w: %v", ErrJSONWrite, err)
	}
	return nil
}
