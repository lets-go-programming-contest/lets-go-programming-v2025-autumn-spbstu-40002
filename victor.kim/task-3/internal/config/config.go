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

func ParseFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}

	defer f.Close()

	return Parse(f)
}

func Parse(r io.Reader) (*Config, error) {
	cfg := new(Config)

	if err := yaml.NewDecoder(r).Decode(cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	return cfg, nil
}
