package currency

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	filePerm = 0o644
	dirPerm  = 0o755
)

var ErrNilCatalog = errors.New("currency catalog is nil")

func OrderByExchange(catalog *CurrencyIndex) error {
	if catalog == nil {
		return fmt.Errorf("%w", ErrNilCatalog)
	}

	sort.Slice(catalog.Currencies, func(i, j int) bool {
		return catalog.Currencies[i].Value > catalog.Currencies[j].Value
	})

	return nil
}

func ExportToJSON(filePath string, catalog *CurrencyIndex) error {
	if catalog == nil {
		return fmt.Errorf("%w", ErrNilCatalog)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(catalog.Currencies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, filePerm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}