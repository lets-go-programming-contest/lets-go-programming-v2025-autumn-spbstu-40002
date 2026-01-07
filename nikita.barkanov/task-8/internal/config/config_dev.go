//go:build dev
// +build dev

package config

import (
	_ "embed"
	"errors"

	"gopkg.in/yaml.v3"
)

var ErrFailUnmarshal = errors.New("yaml error")

//go:embed dev.yaml
var devConfigData []byte

func Load() (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(devConfigData, &cfg); err != nil {
		return nil, ErrFailUnmarshal
	}

	return &cfg, nil
}
