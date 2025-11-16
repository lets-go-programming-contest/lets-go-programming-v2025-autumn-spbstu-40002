package marshaljson

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Svyatoslav2324/task-3/internal/structures"
)

func MarshalJSON(outputFile *os.File, valutes *structures.ValuteStruct) error {
	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", " ")

	err := encoder.Encode(valutes.ValCurs)
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
