//go:build !dev

package config

import (
	_ "embed"
)

//go:embed prod.yaml
var prodConfig []byte

// Load returns production configuration
func LoadConfig() (*Config, error) {
	return loadConfig(prodConfig)
}
