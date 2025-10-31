package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func writeJSON(out []OutputCurrency, path string) error {
	sort.Slice(out, func(i, j int) bool {
		return out[i].Value > out[j].Value
	})
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, dirPerm); err != nil {
				return fmt.Errorf("%w: %v", ErrDirCreate, err)
			}
		}
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, filePerm)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrJSONWrite, err)
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	if err := enc.Encode(out); err != nil {
		f.Close()
		return fmt.Errorf("%w: %v", ErrJSONWrite, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("%w: %v", ErrJSONWrite, err)
	}
	return nil
}
