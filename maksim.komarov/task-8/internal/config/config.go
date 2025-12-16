package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(configYAML, &cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}
