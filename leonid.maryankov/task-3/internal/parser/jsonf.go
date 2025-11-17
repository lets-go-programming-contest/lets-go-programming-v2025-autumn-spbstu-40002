package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	dirPerm  = 0o755
	filePerm = 0o600
)

func SortValute(vs []Valute) {
	sort.Slice(vs, func(i, j int) bool {
		return vs[i].Value > vs[j].Value
	})
}

func SaveToJSON(path string, valutes []Valute) error {
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return fmt.Errorf("error creating directory for %s: %w", path, err)
	}

	data, err := json.MarshalIndent(valutes, "", "  ")
	if err != nil {
		return fmt.Errorf("error in json serialization: %w", err)
	}

	if err := os.WriteFile(path, data, filePerm); err != nil {
		return fmt.Errorf("file write error: %s: %w", path, err)
	}

	return nil
}
