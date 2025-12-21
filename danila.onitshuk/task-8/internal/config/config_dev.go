//go:build dev

package config

import _ "embed"

//go:embed settings/dev.yaml
var configFile []byte
