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
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, fmt.Errorf("%w: %v", ErrConfigParse, err)
	}
	return cfg, nil
}
