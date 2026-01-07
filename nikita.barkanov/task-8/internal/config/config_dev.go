//go:build dev
// +build dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed dev.yaml
var devConfigData []byte

func Load() (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(devConfigData, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
