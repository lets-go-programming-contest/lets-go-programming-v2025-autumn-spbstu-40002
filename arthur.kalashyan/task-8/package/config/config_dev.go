//go:build dev

package config

import _ "embed"

//go:embed dev.yaml
var rawConfig []byte
