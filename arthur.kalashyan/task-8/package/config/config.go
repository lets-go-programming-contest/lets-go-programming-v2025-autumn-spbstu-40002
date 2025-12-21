package config

import "gopkg.in/yaml.v3"

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() Config {
	var cfg Config
	if err := yaml.Unmarshal(rawConfig, &cfg); err != nil {
		panic(err)
	}

	return cfg
}
