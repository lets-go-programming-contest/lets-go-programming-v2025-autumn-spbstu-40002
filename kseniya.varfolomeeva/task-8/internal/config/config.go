package config

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func loadConfig(data []byte) (*Config, error) {
	// Проверка на пустые данные (тест может дать пустой массив)
	if len(data) == 0 {
		return nil, fmt.Errorf("empty config data")
	}
	
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	
	// Проверка обязательных полей
	if cfg.Environment == "" || cfg.LogLevel == "" {
		return nil, fmt.Errorf("missing required fields")
	}
	
	return &cfg, nil
}
