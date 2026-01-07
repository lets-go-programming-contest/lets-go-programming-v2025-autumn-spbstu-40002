//go:build dev
// +build dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed dev.yaml
var devConfig []byte

func load() (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(devConfig, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
