//go:build !dev && !prod

package config

import (
	_ "embed"
)

//go:embed default.yaml
var defaultConfig []byte

func LoadConfig() (*Config, error) {
	return loadConfig(defaultConfig)
}
