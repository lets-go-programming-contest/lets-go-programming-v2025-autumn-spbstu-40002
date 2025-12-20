package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(yamlData, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return &cfg, nil
}
