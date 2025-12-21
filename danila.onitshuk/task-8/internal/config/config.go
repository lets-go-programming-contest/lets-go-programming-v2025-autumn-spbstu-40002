package config

import "gopkg.in/yaml.v3"

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(configFile, &cfg); err != nil {
		return nil, errUnmarshalFailed
	}

	return &cfg, nil
}
