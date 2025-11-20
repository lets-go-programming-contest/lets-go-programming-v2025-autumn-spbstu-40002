package config

import (
	"os"

	"github.com/F0LY/task-3/internal/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(configPath string) *Config {
	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(errors.ErrConfigFileRead.Error() + ": " + err.Error())
	}

	config := &Config{
		InputFile:  "",
		OutputFile: "",
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		panic(errors.ErrConfigFileRead.Error() + ": " + err.Error())
	}

	if config.InputFile == "" || config.OutputFile == "" {
		panic(errors.ErrConfigFieldsMissing.Error())
	}

	return config
}
