//go:build !dev

package config

import _ "embed"

//go:embed settings/prod.yaml
var configFile []byte
