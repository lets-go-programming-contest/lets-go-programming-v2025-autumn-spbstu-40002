package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func getConfig(dataYaml []byte) (*Config, error) {
	var cnfg Config

	err := yaml.Unmarshal(dataYaml, &cnfg)
	if err != nil {
		return nil, fmt.Errorf("error in parsing config: %w", err)
	}

	return &cnfg, nil
}

func (c *Config) PrintToConfig() {
	fmt.Printf("%s %s", c.Environment, c.LogLevel)
}
