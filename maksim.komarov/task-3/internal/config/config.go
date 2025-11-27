package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrOpenConfig       = errors.New("open config")
	ErrDecodeConfigYAML = errors.New("decode config yaml")
	ErrEmptyConfigFlag  = errors.New("config flag is empty")
)

type AppConfig struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(path string) (AppConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return AppConfig{}, fmt.Errorf("%w: %s", ErrOpenConfig, err.Error())
	}

	defer func() { _ = file.Close() }()

	var cfg AppConfig

	dec := yaml.NewDecoder(file)

	if err := dec.Decode(&cfg); err != nil {
		return AppConfig{}, fmt.Errorf("%w: %s", ErrDecodeConfigYAML, err.Error())
	}

	return cfg, nil
}
