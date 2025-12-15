package utils

import (
	"fmt"
	"os"
)

func Execute(configPath string) error {
	cfg, err := loadConfig(configPath)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrXMLRead.Error(), err)
	}

	root, err := parseXML(data)
	if err != nil {
		return err
	}

	out := convertValutes(root)

	if err = writeJSON(out, cfg.OutputFile); err != nil {
		return err
	}

	return nil
}
