package output

import (
	"encoding/json"
	"os"
	"github.com/bolatbyek/task-3/internal/cbr"
)

// Writer handles output operations
type Writer struct{}

// NewWriter creates a new output writer
func NewWriter() *Writer {
	return &Writer{}
}

// SaveToJSON saves currencies to JSON file
func (w *Writer) SaveToJSON(currencies []cbr.Currency, filename string) error {
	// Create directory if it doesn't exist
	if len(filename) > 0 {
		lastSlash := -1
		for i := len(filename) - 1; i >= 0; i-- {
			if filename[i] == '/' || filename[i] == '\\' {
				lastSlash = i
				break
			}
		}
		if lastSlash >= 0 {
			dir := filename[:lastSlash]
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				panic("Failed to create directory: " + err.Error())
			}
		}
	}

	// Create output file
	file, err := os.Create(filename)
	if err != nil {
		panic("Failed to create output file: " + err.Error())
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(currencies)
	if err != nil {
		panic("Failed to encode JSON: " + err.Error())
	}

	return nil
}