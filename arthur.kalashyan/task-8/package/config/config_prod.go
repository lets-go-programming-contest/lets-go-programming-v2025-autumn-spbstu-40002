//go:build !dev

package config

import _ "embed"

//go:embed prod.yaml
var rawConfig []byte
