package config

import (
    "fmt"
    "gopkg.in/yaml.v3"
)

type Config struct {
    Environment string `yaml:"environment"`
    LogLevel    string `yaml:"log_level"`
}

func loadConfig(data []byte) (*Config, error) {
    var cfg Config
    err := yaml.Unmarshal(data, &cfg)
    if err != nil {
        return nil, fmt.Errorf("error in unmarshal file: %w", err)
    }
    return &cfg, nil
}
