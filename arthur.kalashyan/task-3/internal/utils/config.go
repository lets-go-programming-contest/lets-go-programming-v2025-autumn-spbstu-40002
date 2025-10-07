package utils

import (
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		_ = file.Close()
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	if strings.TrimSpace(cfg.InputFile) == "" || strings.TrimSpace(cfg.OutputFile) == "" {
		panic("invalid config")
	}

	return cfg
}
