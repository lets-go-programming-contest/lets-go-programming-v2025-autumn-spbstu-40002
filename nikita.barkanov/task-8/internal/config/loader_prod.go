//go:build !dev
// +build !dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var prodConfigData []byte

func Load() (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(prodConfigData, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
