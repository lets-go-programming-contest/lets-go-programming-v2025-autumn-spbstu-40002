package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func GetConfig() (*Config, error) {
	return nil, fmt.Errorf("config: use build tags 'dev' or default 'prod'")
}

func loadConfig(data []byte) (*Config, error) {
	var cfg Config
	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) PrintConfig() {
	fmt.Printf("%s %s\n", c.Environment, c.LogLevel)
}
