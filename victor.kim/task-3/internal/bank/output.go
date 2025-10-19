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
			return nil, fmt.Errorf("invalid type of value: %w", parseErr)
		}

		// NOTE: Do NOT divide by Nominal — tests expect the raw parsed value.
		out = append(out, outputCurrency{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    parsed,
		})
	}

	return out, nil
}

func (b outputBank) sortByValueDesc() {
	sort.Slice(b, func(i, j int) bool {
		return b[i].Value > b[j].Value
	})
}

func (bank *Bank) EncodeJSON(writer io.Writer) error {
	out, err := buildOutput(bank)
	if err != nil {
		return err
	}

	out.sortByValueDesc()

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(out); err != nil {
		return fmt.Errorf("encoding bank: %w", err)
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

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer func() { _ = file.Close() }()

	return b.EncodeJSON(file)
}
