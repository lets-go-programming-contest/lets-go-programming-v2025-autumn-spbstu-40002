//go:build !dev
// +build !dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var prodConfig []byte

func load() (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(prodConfig, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
