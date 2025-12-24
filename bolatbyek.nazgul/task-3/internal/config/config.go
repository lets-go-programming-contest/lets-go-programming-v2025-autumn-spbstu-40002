package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(configPath string) (*Config, error) {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func EnsureOutputDir(outputFile string) error {
	outputDir := filepath.Dir(outputFile)
	return os.MkdirAll(outputDir, 0755)
}
