package config

import (
	"gopkg.in/yaml.v3"
	"os"

	merr "github.com/slendycs/go-lab-3/internal/myerrors"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func ReadConfig(path string, config *Config) {
	// Read data from file.
	data, err := os.ReadFile(path)
	if err != nil {
		panic(merr.ErrNoConfigFileFound)
	}

	// Deserialize data from config file
	err = yaml.Unmarshal(data, config)
	if err != nil {
		panic(merr.ErrFailedToDeserializeConfig)
	}
}
