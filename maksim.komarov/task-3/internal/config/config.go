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

	file, err := os.Open(path)
	if err != nil {
		return AppConfig{}, fmt.Errorf("%s: %w", ErrOpenConfig, err)
	}

	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&cfg); err != nil {
		return AppConfig{}, fmt.Errorf("%s: %w", ErrDecodeConfig, err)
	}

	return cfg, nil
}
