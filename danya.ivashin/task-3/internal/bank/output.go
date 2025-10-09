package bank

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Danya-byte/task-3/pkg/must"
)

type outputCurrency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

type outputBank []outputCurrency

func fetchOutput(b *Bank) (outputBank, error) {
	out := make(outputBank, len(b.Currencies))

	for index, currency := range b.Currencies {
		val, err := strconv.ParseFloat(strings.Replace(currency.Value, ",", ".", 1), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid type of value: %w", err)
		}

		out[index] = outputCurrency{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    val,
		}
	}

	return out, nil
}

func (b outputBank) sortByValueDown() {
	sort.Slice(
		b,
		func(i, j int) bool {
			return b[i].Value > b[j].Value
		},
	)
}

func (b *Bank) EncodeJSON(writer io.Writer) error {
	out, err := fetchOutput(b)
	if err != nil {
		return err
	}

	out.sortByValueDown()

	encoder := json.NewEncoder(writer)

	encoder.SetIndent("", "  ")

	if err := encoder.Encode(&out); err != nil {
		return fmt.Errorf("encoding bank: %w", err)
	}

	return nil
}

func (b *Bank) EncodeFileJSON(path string) error {
	const permissions = 0o755

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, permissions); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer must.Close(path, file)

	return b.EncodeJSON(file)
}
