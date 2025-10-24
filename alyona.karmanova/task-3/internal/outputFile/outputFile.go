package outputfile

import (
	"fmt"
	"os"
	"path/filepath"
)

func EnsureOutputDir(path string) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("couldn't create a directory %s: %w", dir, err)
	}

	return nil
}
