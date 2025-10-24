package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/megurumacabre/task-3/internal/convert"
)

const (
	permDir  = 0o755
	permFile = 0o644
)

var (
	ErrMakeOutputDir    = errors.New("make output dir")
	ErrCreateOutputFile = errors.New("create output file")
	ErrWriteJSON        = errors.New("write output json")
)

func WriteJSON(path string, data []convert.CurrencyOut) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, permDir); err != nil {
			return fmt.Errorf("%w: %v", ErrMakeOutputDir, err)
		}
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, permFile)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCreateOutputFile, err)
	}
	defer f.Close()

	blob, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWriteJSON, err)
	}

	if _, err := f.Write(append(blob, '\n')); err != nil {
		return fmt.Errorf("%w: %v", ErrWriteJSON, err)
	}

	return nil
}
