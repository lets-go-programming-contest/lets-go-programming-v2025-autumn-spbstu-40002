package writer

import (
	"encoding/json"
	"os"

	"github.com/XShaygaND/task-3/internal/parser"
)

func WriteJSON(filePath string, data []parser.Item) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
