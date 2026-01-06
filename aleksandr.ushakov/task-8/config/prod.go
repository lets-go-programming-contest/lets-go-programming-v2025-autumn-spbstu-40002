//go:build !dev

package config

import (
	"embed"
)

//go:embed prod.yaml
var conf_prod embed.FS

func init() {
	Cnf = *ParseConfig(&conf_prod, "prod.yaml")
}
