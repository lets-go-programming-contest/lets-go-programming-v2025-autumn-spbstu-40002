package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

const errUnmarshalConfig = "failed to unmarshal config"

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (*Config, error) {
	var cfg Config
	data := getConfigData()

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalConfig, err)
	}
	return &cfg, nil
}
