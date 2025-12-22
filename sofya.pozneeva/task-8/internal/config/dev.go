//go:build dev

package config

import _ "embed"

//go:embed dev.yaml
var config []byte
