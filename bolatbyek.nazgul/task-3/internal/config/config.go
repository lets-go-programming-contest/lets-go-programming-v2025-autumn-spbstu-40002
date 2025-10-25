package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

// LoadConfig loads configuration from YAML file
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		panic("Failed to read config file: " + err.Error())
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic("Failed to parse config file: " + err.Error())
	}

	return &config, nil
}