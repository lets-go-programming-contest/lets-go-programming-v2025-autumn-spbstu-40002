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

func buildOutput(bankData *Bank) (outputBank, error) {
	result := make(outputBank, 0, len(bankData.Currencies))

	for _, curr := range bankData.Currencies {
		raw := strings.TrimSpace(curr.Value)
		raw = strings.Replace(raw, ",", ".", 1)

		parsedVal, parseErr := strconv.ParseFloat(raw, 64)
		if parseErr != nil {
			return nil, fmt.Errorf("parse value failed: %w", parseErr)
		}

		result = append(result, outputCurrency{
			NumCode:  curr.NumCode,
			CharCode: curr.CharCode,
			Value:    parsedVal,
		})
	}

	return result, nil
}

func (outBank outputBank) sortByValueDesc() {
	sort.Slice(outBank, func(i, j int) bool {
		return outBank[i].Value > outBank[j].Value
	})
}

func (bankData *Bank) WriteJSON(writer io.Writer) error {
	out, err := buildOutput(bankData)
	if err != nil {
		return err
	}

	out.sortByValueDesc()

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(out); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func (bankData *Bank) WriteJSONFile(path string) error {
	const perm = 0o755

	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, perm); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	fileHandle, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	defer func() { _ = fileHandle.Close() }()

	return bankData.WriteJSON(fileHandle)
}
