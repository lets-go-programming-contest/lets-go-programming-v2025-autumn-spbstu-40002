package jsonutils

import (
	stdjson "encoding/json"
	"os"
	"path/filepath"
	"sort"

	"github.com/xkoex/task-3/internal/xmlutils"
)

const dirPerm = 0o755

func SortCurrencies(currencies []xmlutils.Currency) {
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].FloatValue() > currencies[j].FloatValue()
	})
}

func WriteJSON(currencies []xmlutils.Currency, path string) {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, dirPerm); err != nil {
		panic(err)
	}

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	type Out struct {
		NumCode  int     `json:"num_code"`
		CharCode string  `json:"char_code"`
		Value    float64 `json:"value"`
	}

	result := make([]Out, 0, len(currencies))
	for _, c := range currencies {
		result = append(result, Out{
			NumCode:  c.NumCode,
			CharCode: c.CharCode,
			Value:    c.FloatValue(),
		})
	}

	encoder := stdjson.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(result); err != nil {
		panic(err)
	}
}
