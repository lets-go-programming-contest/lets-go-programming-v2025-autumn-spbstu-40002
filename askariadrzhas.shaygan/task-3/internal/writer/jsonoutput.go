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

	file, err := os.Create(outputPath)
	if err != nil {
		return
	}

	defer func() {
		_ = file.Close()
	}()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "    ")

	if err := enc.Encode(data); err != nil {
		panic("cannot write JSON data: " + err.Error())
	}
}
