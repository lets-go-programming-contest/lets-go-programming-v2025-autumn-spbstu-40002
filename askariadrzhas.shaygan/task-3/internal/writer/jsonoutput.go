package writer

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const dirPerm = 0o755

func SaveAsJSON(data interface{}, outputPath string) {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		panic("cannot create output directory: " + err.Error())
	}

	f, err := os.Create(outputPath)
	if err != nil {
		panic("cannot create output file: " + err.Error())
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			panic("cannot close file: " + cerr.Error())
		}
	}()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")

	if err := enc.Encode(data); err != nil {
		panic("cannot write JSON data: " + err.Error())
	}
}
