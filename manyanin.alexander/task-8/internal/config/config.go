package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func LoadAppConfig() (*AppConfig, error) {
	var cfg AppConfig

	err := yaml.Unmarshal(configFileData, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
