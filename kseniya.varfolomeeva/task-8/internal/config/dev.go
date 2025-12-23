//go:build dev

package config

import (
	_ "embed"
)

//go:embed dev.yaml
var devConfig []byte

func LoadConfig() (*Config, error) {
	return loadConfig(developmentConfig)
}
