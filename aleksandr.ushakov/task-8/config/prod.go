//go:build !dev

package config

import (
	"embed"

	"github.com/rachguta/task-8/myerrors"
	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var confFile embed.FS

func GetConfig() *Config {
	// read config
	data, err := confFile.ReadFile("prod.yaml")
	if err != nil {
		panic(myerrors.ErrConfigRead)
	}
	// parse config
	var cnf Config

	err = yaml.Unmarshal(data, &cnf)
	if err != nil || cnf.Environment == "" || cnf.Loglevel == "" {
		panic(myerrors.ErrConfigParse)
	}

	return &cnf
}
