package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/manyanin.alexander/task-3/internal/errors"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(configPath string) *Config {
	if configPath == "" {
		panic(errors.ErrConfigPathEmpty)
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		panic(errors.ErrConfigRead.Error() + ": " + configPath)
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		panic(errors.ErrConfigParse.Error() + ": " + err.Error())
	}

	if config.InputFile == "" || config.OutputFile == "" {
		panic(errors.ErrConfigInvalid)
	}

	return &config
}
