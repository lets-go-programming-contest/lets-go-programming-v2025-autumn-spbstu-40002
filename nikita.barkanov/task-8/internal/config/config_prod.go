//go:build !dev
// +build !dev

package config

import (
	_ "embed"
	"errors"

	"gopkg.in/yaml.v3"
)

var ErrFailUnmarshal = errors.New("yaml error")

//go:embed prod.yaml
var prodConfigData []byte

func Load() (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(prodConfigData, &cfg); err != nil {
		return nil, ErrFailUnmarshal
	}

	return &cfg, nil
}
