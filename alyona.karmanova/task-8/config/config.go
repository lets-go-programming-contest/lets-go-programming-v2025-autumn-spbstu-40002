package config

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var ErrUnmarshalConfig = errors.New("failed to unmarshal config")

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (*Config, error) {
	var cfg Config

	data := getConfigData()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnmarshalConfig, err)
	}

	return &cfg, nil
}
