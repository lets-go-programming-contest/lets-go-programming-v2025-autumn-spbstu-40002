//go:build dev
// +build dev

package config

import _ "embed"

//go:embed dev.yaml
var configYAML []byte
