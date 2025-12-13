package config

import (
	"github.com/slendycs/go-lab-8/internal/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (*Config, error) {
	var tempConfig Config

	err := yaml.Unmarshal(configFile, &tempConfig)
	if err != nil {
		return nil, errors.ErrUnmarshalFailed
	}

	return &tempConfig, nil
}
