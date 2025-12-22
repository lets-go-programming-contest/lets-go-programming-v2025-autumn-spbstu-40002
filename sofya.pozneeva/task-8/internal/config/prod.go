//go:build !dev || prod

package config

import _ "embed"

//go:embed prod.yaml
var config []byte
