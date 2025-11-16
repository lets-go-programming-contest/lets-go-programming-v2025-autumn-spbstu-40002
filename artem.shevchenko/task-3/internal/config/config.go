package config

import (
	"os"

	merr "github.com/slendycs/go-lab-3/internal/myerrors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func ReadConfig(path string, config *Config) error {
	// Read data from file.
	data, err := os.ReadFile(path)
	if err != nil {
		return merr.ErrNoConfigFileFound
	}

	// Deserialize data from config file
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return merr.ErrFailedToDeserializeConfig
	}

	return nil
}
