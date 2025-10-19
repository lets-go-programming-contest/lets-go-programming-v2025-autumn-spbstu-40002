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
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}

	defer func() { _ = file.Close() }()

	return Parse(file)
}

func Parse(r io.Reader) (*Config, error) {
	cfg := new(Config)

	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("decoding config file: %w", err)
	}

	return cfg, nil
}
