package bank

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/Danya-byte/task-3/pkg/must"
)

func (b *Bank) SortByValueDesc() {
	sort.Slice(b.Currencies, func(i, j int) bool {
		return b.Currencies[i].Value > b.Currencies[j].Value
	})
}

func (b *Bank) EncodeJSON(w io.Writer) error {
	b.SortByValueDesc()

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	if err := enc.Encode(b); err != nil {
		return fmt.Errorf("encoding bank: %w", err)
	}
	return nil
}

func (b *Bank) EncodeFileJSON(path string) error {
	const perm = 0o755

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer must.Close(path, f)

	return b.EncodeJSON(f)
}
