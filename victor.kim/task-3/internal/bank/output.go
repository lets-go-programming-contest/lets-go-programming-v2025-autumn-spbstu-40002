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

	"github.com/victor.kim/task-3/pkg/must"
)

type outputCurrency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

type outputBank []outputCurrency

func buildOutput(b *Bank) (outputBank, error) {
	out := make(outputBank, 0, len(b.Currencies))

	for _, currency := range b.Currencies {
		raw := strings.TrimSpace(currency.Value)
		raw = strings.Replace(raw, ",", ".", 1)

		parsed, parseErr := strconv.ParseFloat(raw, 64)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid type of value for %q: %w", currency.CharCode, parseErr)
		}

		out = append(out, outputCurrency{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    parsed,
		})
	}

	return out, nil
}

func (ob outputBank) sortByValueDesc() {
	sort.Slice(ob, func(i, j int) bool {
		return ob[i].Value > ob[j].Value
	})
}

func (b *Bank) EncodeJSON(writer io.Writer) error {
	out, err := buildOutput(b)
	if err != nil {
		return err
	}

	out.sortByValueDesc()

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")

	if encErr := encoder.Encode(out); encErr != nil {
		return fmt.Errorf("encoding bank: %w", encErr)
	}

	return nil
}

func (b *Bank) EncodeFileJSON(path string) error {
	const perm = 0o755

	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, perm); err != nil {
			return fmt.Errorf("create dir: %w", err)
		}
	}

	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer must.Close(path, outFile)

	return b.EncodeJSON(outFile)
}
