package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	InputPath  string `yaml:"input"`
	OutputPath string `yaml:"output"`
}

func LoadSettings() (*Settings, error) {
	configPath := "config.yaml"

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Settings
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.InputPath == "" || cfg.OutputPath == "" {
		return nil, errors.New("missing required paths in configuration")
	}

	return &cfg, nil
}
