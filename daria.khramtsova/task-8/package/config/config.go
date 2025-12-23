package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(conf, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling error: %w", err)
	}

	return &cfg, nil
}
