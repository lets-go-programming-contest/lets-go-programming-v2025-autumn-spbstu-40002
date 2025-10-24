package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	errOpeningConfigFile = errors.New("error occurred while opening config file")
	errParsingYAML       = errors.New("error occurred while parsing yaml")
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func GetConfig(path string) (Config, error) {
	var cfg Config

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, errOpeningConfigFile
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, errParsingYAML
	}

	return cfg, nil
}
