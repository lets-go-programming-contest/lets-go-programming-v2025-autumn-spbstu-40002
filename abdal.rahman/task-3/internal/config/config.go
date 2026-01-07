package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func NewConfig(path string, cfg *Config) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(ErrPath)
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil || cfg.InputFile == "" || cfg.OutputFile == "" {
		panic(ErrCfg)
	}
}
