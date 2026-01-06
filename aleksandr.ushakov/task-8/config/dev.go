//go:build dev

package config

import (
	"embed"
)

//go:embed dev.yaml
var conf_dev embed.FS

func init() {
	Cnf = *ParseConfig(&conf_dev, "dev.yaml")
}
