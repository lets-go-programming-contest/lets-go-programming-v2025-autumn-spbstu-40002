package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/megurumacabre/task-3/internal/cbr"
)

var (
	ErrMakeOutputDir    = errors.New("make output dir")
	ErrCreateOutputFile = errors.New("create output file")
	ErrWriteJSON        = errors.New("write json")
)

const dirPerm = 0o755

func WriteJSON(path string, data []cbr.Currency) error {
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return fmt.Errorf("%w: %s", ErrMakeOutputDir, err.Error())
	}

	f, err := os.Create(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCreateOutputFile, err.Error())
	}
	defer func() { _ = f.Close() }()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("%w: %s", ErrWriteJSON, err.Error())
	}

	return nil
}
