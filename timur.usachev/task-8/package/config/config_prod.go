//go:build !dev

package config

import (
	_ "embed"
	"log"
)

//go:embed prod.yaml
var prodConfig []byte

func Load() *Config {
	cfg, err := parse(prodConfig)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
