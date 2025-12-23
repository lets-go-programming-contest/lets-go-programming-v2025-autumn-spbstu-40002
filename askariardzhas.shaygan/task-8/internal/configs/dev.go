//go:build dev

package configs

import _ "embed"

//go:embed config/dev.yaml
var configData []byte
