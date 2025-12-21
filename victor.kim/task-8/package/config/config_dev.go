//go:build dev

package config

import _ "embed"

//go:embed dev.yaml
var ConfigFile []byte
