package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrMakeOutputDir    = errors.New("make output dir")
	ErrCreateOutputFile = errors.New("create output file")
	ErrWriteJSON        = errors.New("write json")
)

func WriteJSON(path string, data any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("%s: %w", ErrMakeOutputDir, err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrCreateOutputFile, err)
	}
	defer func() { _ = f.Close() }()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("%s: %w", ErrWriteJSON, err)
	}

	return nil
}
