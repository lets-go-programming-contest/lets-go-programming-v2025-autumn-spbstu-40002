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

var (
	errMissingPaths = errors.New("missing required paths in configuration")
)

func LoadSettings() (*Settings, error) {
	cfg := &Settings{}

	if data, err := os.ReadFile("config.yaml"); err == nil {
		// file exists, parse YAML
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config: %w", err)
		}
	} else {
		cfg.InputPath = os.Getenv("LGP_TASK_INPUT_PATH")
		cfg.OutputPath = os.Getenv("LGP_TASK_OUTPUT_PATH")
	}

	if cfg.InputPath == "" || cfg.OutputPath == "" {
		return nil, fmt.Errorf("invalid configuration: %w", errMissingPaths)
	}

	return cfg, nil
}
