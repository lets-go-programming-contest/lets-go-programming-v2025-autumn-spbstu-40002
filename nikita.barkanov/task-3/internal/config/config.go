package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrConfigFileOpen    = errors.New("failed to open config file")
	ErrYAMLParsing       = errors.New("failed to parse YAML config")
	ErrInputFileMissing  = errors.New("inputFile is required in config")
	ErrOutputFileMissing = errors.New("outputFile is required in config")
)

type Config struct {
	InputFile  string `yaml:"inputFile"`
	OutputFile string `yaml:"outputFile"`
}

func GetConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrConfigFileOpen, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrYAMLParsing, err)
	}

	if cfg.InputFile == "" {
		return nil, ErrInputFileMissing
	}
	if cfg.OutputFile == "" {
		return nil, ErrOutputFileMissing
	}

	return &cfg, nil
}
