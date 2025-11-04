package writer

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func SaveAsJSON(data interface{}, outputPath string) {
	dir := filepath.Dir(outputPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic("cannot create output directory: " + err.Error())
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic("cannot create output file: " + err.Error())
	}
	defer outputFile.Close()

	jsonEncoder := json.NewEncoder(outputFile)
	jsonEncoder.SetIndent("", "    ")

	err = jsonEncoder.Encode(data)
	if err != nil {
		panic("cannot write JSON data: " + err.Error())
	}
}
