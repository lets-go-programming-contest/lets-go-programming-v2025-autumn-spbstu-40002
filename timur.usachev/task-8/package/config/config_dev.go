//go:build dev

package config

import (
	_ "embed"
	"log"
)

//go:embed dev.yaml
var devConfig []byte

func Load() *Config {
	cfg, err := parse(devConfig)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
