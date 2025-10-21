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

	defer func() { _ = fileHandle.Close() }()

	return DecodeConfig(fileHandle)
}

func DecodeConfig(r io.Reader) (*Config, error) {
	configData := new(Config)

	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(configData); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return configData, nil
}
