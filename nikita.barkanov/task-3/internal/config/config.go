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
	Input  string `yaml:"input-file"`
	Output string `yaml:"output-file"`
}

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfigFileOpen, err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "error: failed to close file %v\n", closeErr)
		}
	}()

	var cfg Config
	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrYAMLParsing, err)
	}

	return &cfg, nil
}
