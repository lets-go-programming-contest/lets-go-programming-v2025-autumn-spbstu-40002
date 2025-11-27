package config

import (
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Input  string `yaml:"input-file"`
	Output string `yaml:"output-file"`
}

func LoadConfigFromFile(path string) (*Config, error) {
	fileHandle, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %w", err)
	}

	loadedConfig, decodeErr := DecodeConfig(fileHandle)
	if closeErr := fileHandle.Close(); decodeErr != nil {
		if closeErr != nil {
			return nil, fmt.Errorf("failed to load config: %w", errors.Join(decodeErr, closeErr))
		}

		return nil, fmt.Errorf("failed to load config: %w", decodeErr)
	} else if closeErr != nil {
		return nil, fmt.Errorf("failed to close config file: %w", closeErr)
	}

	return loadedConfig, nil
}

func DecodeConfig(r io.Reader) (*Config, error) {
	configData := new(Config)

	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(configData); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return configData, nil
}
