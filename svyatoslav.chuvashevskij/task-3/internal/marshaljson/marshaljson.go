package marshaljson

import (
	"encoding/json"
	"os"

	"github.com/Svyatoslav2324/task-3/internal/data"
)

func MarshalJSON(outputFile *os.File, valutes *data.DataStruct) {
	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", " ")
	encoder.Encode(*valutes)
}
