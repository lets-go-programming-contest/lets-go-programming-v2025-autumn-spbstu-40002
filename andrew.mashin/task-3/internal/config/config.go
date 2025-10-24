package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	errOpenYamlFile      = errors.New("error opening config file")
	errUnmarshalYamlFile = errors.New("error unmarshalling yaml file")
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errOpenYamlFile, err)
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errUnmarshalYamlFile, err)
	}

	return &config, nil
}
