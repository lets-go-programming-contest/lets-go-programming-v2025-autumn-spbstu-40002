//go:build prod || !dev

package configs

import _ "embed"

//go:embed config/prod.yaml
var dataYaml []byte
