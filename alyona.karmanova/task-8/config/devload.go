//go:build dev

package config

import _ "embed"

//go:embed dev.yaml
var configData []byte

func getConfigData() []byte {
	return configData
}
