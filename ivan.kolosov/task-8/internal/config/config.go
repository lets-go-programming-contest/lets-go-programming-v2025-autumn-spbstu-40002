package config

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var errUnmarshaling = errors.New("an error occurred while unmarshaling config file")

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (*Config, error) {
	var cfg Config

	err := yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errUnmarshaling, err)
	}

	return &cfg, nil
}
