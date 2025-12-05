package config

import (
	"fmt"
	"os"

	"github.com/manyanin.alexander/task-3/internal/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrConfigRead, err)
	}

	config := &Config{
		InputFile:  "",
		OutputFile: "",
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrConfigParse, err)
	}

	if config.InputFile == "" || config.OutputFile == "" {
		return nil, errors.ErrConfigInvalid
	}

	return config, nil
}
