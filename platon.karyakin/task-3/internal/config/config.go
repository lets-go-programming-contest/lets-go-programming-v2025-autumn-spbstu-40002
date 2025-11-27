package config

import (
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
	
	closeError := fileHandle.Close()
	if decodeErr != nil {
		if closeError != nil {
			return nil, fmt.Errorf("decode config error: %v; close error: %w", decodeErr, closeError)
		}
		return nil, decodeErr
	}
	if closeError != nil {
		return nil, fmt.Errorf("failed to close config file: %w", closeError)
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
