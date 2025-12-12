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
	if configContent == nil {
		return nil, fmt.Errorf("config content is empty")
	}

	var cfg AppConfig
	err := yaml.Unmarshal(configContent, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
