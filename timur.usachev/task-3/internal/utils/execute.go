package utils

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func Execute(configPath string) error {
	cfg, err := loadConfig(configPath)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrXMLRead, err)
	}
	var root valCurs
	if err := xml.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("%w: %v", ErrXMLParse, err)
	}
	out := make([]OutputCurrency, 0, len(root.Valutes))
	for _, v := range root.Valutes {
		nc := strings.TrimSpace(v.NumCode)
		cc := strings.TrimSpace(v.CharCode)
		valStr := strings.TrimSpace(v.Value)
		valStr = strings.ReplaceAll(valStr, ",", ".")
		numCode, err := strconv.Atoi(nc)
		if err != nil {
			return fmt.Errorf("%w: invalid num code %q: %v", ErrXMLParse, nc, err)
		}
		value, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return fmt.Errorf("%w: invalid value %q: %v", ErrXMLParse, valStr, err)
		}
		out = append(out, OutputCurrency{
			NumCode:  numCode,
			CharCode: cc,
			Value:    value,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Value > out[j].Value
	})
	if err := ensureDirForFile(cfg.OutputFile); err != nil {
		return err
	}
	f, err := os.OpenFile(cfg.OutputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, filePerm)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrJSONWrite, err)
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	if err := enc.Encode(out); err != nil {
		f.Close()
		return fmt.Errorf("%w: %v", ErrJSONWrite, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("%w: %v", ErrJSONWrite, err)
	}
	return nil
}

func ensureDirForFile(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir == "." || dir == "" {
		return nil
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return fmt.Errorf("%w: %v", ErrDirCreate, err)
		}
	}
	return nil
}
