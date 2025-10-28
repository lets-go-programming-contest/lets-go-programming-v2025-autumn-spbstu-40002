package marshaljson

import (
	"encoding/json"
	"os"

	"github.com/Svyatoslav2324/task-3/internal/data"
)

func MarshalJSON(outputFile *os.File, valutes *data.DataStruct) error {
	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", " ")

	err := encoder.Encode(valutes.ValCurs)

	if err != nil {

		return err
	}

	return nil
}
