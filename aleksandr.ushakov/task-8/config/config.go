package config

import (
	"embed"

	"github.com/rachguta/task-8/myerrors"
	"gopkg.in/yaml.v3"
)

var Cnf Config

type Config struct {
	Environment string `yaml:"environment"`
	Log_level   string `yaml:"log_level"`
}

func ParseConfig(fs *embed.FS, configPath string) *Config {
	// read config
	data, err := fs.ReadFile(configPath)
	if err != nil {
		panic(myerrors.ErrConfigRead)
	}
	// parse config
	var cnf Config

	err = yaml.Unmarshal(data, &cnf)
	if err != nil || cnf.Environment == "" || cnf.Log_level == "" {
		panic(myerrors.ErrConfigParse)
	}

	return &cnf
}
