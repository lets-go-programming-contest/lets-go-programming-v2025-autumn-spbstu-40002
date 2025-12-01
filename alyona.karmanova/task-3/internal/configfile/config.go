package configfile

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

var (
	errConfigRead    = errors.New("error with read the config file")
	errConfigParsing = errors.New("error with config yaml file parsing")
)

func GetConfigStruct(path string) (Config, error) {
	var cfg Config

	dataConf, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("%w: %w", errConfigRead, err)
	}

	err = yaml.Unmarshal(dataConf, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("%w: %w", errConfigParsing, err)
	}

	return cfg, nil
}
