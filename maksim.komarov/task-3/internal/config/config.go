package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrOpenConfig   = errors.New("open config file")
	ErrDecodeConfig = errors.New("decode config")
)

type AppConfig struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(path string) (AppConfig, error) {
	var cfg AppConfig

	f, err := os.Open(path)
	if err != nil {
		return AppConfig{}, fmt.Errorf("%w: %v", ErrOpenConfig, err)
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return AppConfig{}, fmt.Errorf("%w: %v", ErrDecodeConfig, err)
	}

	return cfg, nil
}
