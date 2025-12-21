package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrConfigPathRequired = errors.New("config path is required")
	ErrFileNotFound       = errors.New("no such file or directory")
	ErrInvalidYAML        = errors.New("did not find expected key")
)

type AppConfig struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	if configPath == "" {
		return nil, ErrConfigPathRequired
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, ErrFileNotFound
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, ErrInvalidYAML
	}

	return &config, nil
}
