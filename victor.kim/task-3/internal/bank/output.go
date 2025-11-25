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

func (b *Bank) EncodeJSON(w io.Writer) error {
	// prepare currencies for output
	prepared := make([]map[string]interface{}, 0, len(b.Currencies))

	for _, cur := range b.Currencies {
		v := strings.Replace(strings.TrimSpace(cur.Value), ",", ".", 1)
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("parse currency value: %w", err)
		}

		prepared = append(prepared, map[string]interface{}{
			"num_code":  cur.NumCode,
			"char_code": cur.CharCode,
			"value":     parsed,
		})
	}

	sort.Slice(prepared, func(i, j int) bool {
		return prepared[i]["value"].(float64) > prepared[j]["value"].(float64)
	})

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	return enc.Encode(prepared)
}

func (b *Bank) EncodeFileJSON(path string) error {
	const perms = 0o755

	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, perms); err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}

	defer must.Close(path, f)

	return b.EncodeJSON(f)
}
