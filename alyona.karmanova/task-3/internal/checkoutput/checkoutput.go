package checkoutput

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var errCreatingDir = errors.New("error with creating output directory")

func EnsureOutputDir(path string) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("%w: %s: %w", errCreatingDir, dir, err)
	}

	return nil
}
