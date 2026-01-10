//go:build !dev

package config

import (
	_ "embed"
)

//go:embed prod.yaml
var prodConfig []byte

func LoadConfig() (*Config, error) {
	return loadConfig(prodConfig)
}
