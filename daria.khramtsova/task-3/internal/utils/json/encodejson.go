package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"github.com/hehemka/task-3/internal/utils/xml"
)

func SortCurrencies(currencies []xml.Currency) {
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].FloatValue() > currencies[j].FloatValue()
	})
}

func WriteJSON(currencies []xml.Currency, path string) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

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

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		panic(err)
	}
}
