//go:build !dev

package config

import (
	_ "embed"
)

//go:embed prod.yaml
var productionConfig []byte

// Load returns production configuration
func Load() (*Config, error) {
	return parseYAML(productionConfig)
}
