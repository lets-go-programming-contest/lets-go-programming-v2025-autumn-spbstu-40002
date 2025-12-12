package config

import (
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
