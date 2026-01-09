package core

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	dirPerm  = 0o755
	filePerm = 0o600
)

func PrepareJSON(valutes *Valutes) ([]byte, error) {
	out := make([]Valute, 0)
	for _, v := range *valutes {
		out = append(out, Valute{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    v.Value,
		})
	}
	data, err := json.MarshalIndent(out, "", "    ")
	if err != nil {
		return nil, err
	}
	return data, nil
}

func JSONWrite(cfg *Config, data []byte) error {
	dir := filepath.Dir(cfg.OutputFile)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return err
	}
	if err := os.WriteFile(cfg.OutputFile, data, filePerm); err != nil {
		return err
	}
	return nil
}
