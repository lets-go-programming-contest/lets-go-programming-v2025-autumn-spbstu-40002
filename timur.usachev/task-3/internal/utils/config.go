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

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("%s: %w", ErrConfigRead.Error(), err)
	}

	var cfg Config

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("%s: %w", ErrConfigParse.Error(), err)
	}

	return cfg, nil
}
