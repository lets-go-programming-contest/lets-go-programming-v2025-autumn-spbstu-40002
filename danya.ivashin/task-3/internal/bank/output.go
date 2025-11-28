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

const Permissions = 0o755

type outputCurrency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

type outputBank []outputCurrency

func fetchOutput(b *Bank) (outputBank, error) {
	output := make(outputBank, len(b.Currencies))

	for index, currency := range b.Currencies {
		valueFloat, err := strconv.ParseFloat(strings.Replace(currency.Value, ",", ".", 1), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid type of value: %w", err)
		}

		output[index] = outputCurrency{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    valueFloat,
		}
	}

	return output, nil
}

func (b outputBank) sortByValueDown() {
	sort.Slice(b, func(firstIndex, secondIndex int) bool {
		return b[firstIndex].Value > b[secondIndex].Value
	})
}

func (b *Bank) EncodeJSON(writer io.Writer) error {
	output, err := fetchOutput(b)
	if err != nil {
		return err
	}

	output.sortByValueDown()

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(&output); err != nil {
		return fmt.Errorf("encoding bank: %w", err)
	}

	return nil
}

func (b *Bank) EncodeFileJSON(path string) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, Permissions); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer must.Close(path, file)

	return b.EncodeJSON(file)
}
