package config

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var errUnmarshal = errors.New("failed to unmarshal yaml")

type Cfg struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func LoadCfg() (*Cfg, error) {
	var cfg Cfg

	err := yaml.Unmarshal(cfgFile, &cfg)
	if err != nil {
		return &Cfg{}, fmt.Errorf("%w: %w", errUnmarshal, err)
	}

	return &cfg, nil
}
