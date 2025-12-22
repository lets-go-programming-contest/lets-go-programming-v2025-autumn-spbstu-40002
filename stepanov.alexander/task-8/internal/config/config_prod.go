//go:build !dev

package config

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var prodYaml []byte

func Load() (Config, error) {
	var cfg Config

	err := yaml.Unmarshal(prodYaml, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal prod config: %w", err)
	}

	return cfg, nil
}
