package writer

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const dirPermissions = 0o755

func SaveAsJSON(data interface{}, outputPath string) {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, dirPermissions); err != nil {
		panic("cannot create output directory: " + err.Error())
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("cannot create output file: " + err.Error())
	}

	defer func() {
		if cerr := outputFile.Close(); cerr != nil {
			panic("cannot close file: " + cerr.Error())
		}
	}()

	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "    ")

	if err = encoder.Encode(data); err != nil {
		panic("cannot write JSON data: " + err.Error())
	}
}
