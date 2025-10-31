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
		return fmt.Errorf("%w", ErrXMLRead)
	}

	root, err := parseXML(data)
	if err != nil {
		return err
	}

	out := convertValutes(root)

	err = writeJSON(out, cfg.OutputFile)
	if err != nil {
		return err
	}

	return nil
}
