package config

import (
	"errors"
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type AppSettings struct {
	SourcePath string `yaml:"inputFile"`
	TargetPath string `yaml:"outputFile"`
}

func LoadSettings() (*AppSettings, error) {
	configPath := flag.String("config", "", "Configuration file path")
	flag.Parse()

	if *configPath == "" {
		return nil, errors.New("no config path provided")
	}

	data, err := os.ReadFile(*configPath)
	if err != nil {
		return nil, err
	}

	var cfg AppSettings
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.SourcePath == "" || cfg.TargetPath == "" {
		return nil, errors.New("missing required paths in configuration")
	}

	return &cfg, nil
}
