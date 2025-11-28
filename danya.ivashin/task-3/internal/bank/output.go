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

const FilePermissions = 0o755

type outputBank []Currency

func fetchOutput(b *Bank) (outputBank, error) {
	out := make(outputBank, len(b.Currencies))

	for i, c := range b.Currencies {
		v, err := strconv.ParseFloat(strings.Replace(c.Value, ",", ".", 1), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid type of value: %w", err)
		}

		out[i] = c
		out[i].Value = fmt.Sprintf("%.4f", v)
	}

	return out, nil
}

func (b outputBank) sortByValueDown() {
	sort.Slice(b, func(i, j int) bool {
		vi, _ := strconv.ParseFloat(b[i].Value, 64)
		vj, _ := strconv.ParseFloat(b[j].Value, 64)
		return vi > vj
	})
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
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, FilePermissions); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer must.Close(path, file)

	return b.EncodeJSON(file)
}
