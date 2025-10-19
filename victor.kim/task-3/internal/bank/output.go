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

	for _, c := range b.Currencies {
		valStr := strings.TrimSpace(c.Value)
		valStr = strings.ReplaceAll(valStr, ",", ".")
		parsed, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for %s: %w", c.CharCode, err)
		}

		nominal := c.Nominal
		if nominal <= 0 {
			nominal = 1
		}

		out = append(out, outputCurrency{
			NumCode:  c.NumCode,
			CharCode: c.CharCode,
			Value:    parsed / float64(nominal),
		})
	}

	return out, nil
}

func (b outputBank) sortByValueDesc() {
	sort.Slice(b, func(i, j int) bool {
		return b[i].Value > b[j].Value
	})
}

func (b *Bank) EncodeJSON(w io.Writer) error {
	out, err := buildOutput(b)
	if err != nil {
		return err
	}

	out.sortByValueDesc()

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	if err := enc.Encode(out); err != nil {
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

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer func() { _ = f.Close() }()

	return b.EncodeJSON(f)
}
