//go:build dev

package config

import _ "embed"

//go:embed dev.yaml
var dataYAML []byte

func Initial() (*Config, error) {
	return getConfig(dataYAML)
}
