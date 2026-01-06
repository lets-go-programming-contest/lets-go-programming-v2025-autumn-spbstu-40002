package config

import (
	"github.com/rachguta/task-8/myerrors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	Loglevel    string `yaml:"log_level"`
}

func GetConfig() *Config {
	// read config
	data, err := confFile.ReadFile(getFilePath())
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
