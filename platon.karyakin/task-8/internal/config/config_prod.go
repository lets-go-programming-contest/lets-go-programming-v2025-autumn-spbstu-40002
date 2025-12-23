//go:build !dev

package config

import _ "embed"

//go:embed prod.yaml
var dataYAML []byte

func Initial() (*Config, error) {
	return getConfig(dataYAML)
}
