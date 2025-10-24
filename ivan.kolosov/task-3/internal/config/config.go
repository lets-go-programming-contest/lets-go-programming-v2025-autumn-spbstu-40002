package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	errOpeningConfigFile = "error occurred while opening config file"
	errParsingYAML       = "error occurred while parsing yaml"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func GetConfig(path string) (Config, error) {
	var cfg Config

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("%s: %w", errOpeningConfigFile, err)
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("%s: %w", errParsingYAML, err)
	}

	return cfg, nil
}
