package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func loadConfig(path string) (Config, error) {
	if path == "" {
		path = defaultConfigPath
	}

	fileData, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("%w", ErrConfigRead)
	}

	var cfg Config

	err = yaml.Unmarshal(fileData, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("%w", ErrConfigParse)
	}

	return cfg, nil
}
