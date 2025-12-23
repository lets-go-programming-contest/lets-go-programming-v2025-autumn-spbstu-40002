//go:build dev

package config

import (
	_ "embed"
)

//go:embed dev.yaml
var developmentConfig []byte

func Load() (*Config, error) {
	return parseYAML(developmentConfig)
}
