//go:build !dev
// +build !dev

package config

import _ "embed"

//go:embed prod.yaml
var configYAML []byte
