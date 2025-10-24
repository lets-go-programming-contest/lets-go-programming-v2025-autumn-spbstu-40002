package configfile

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func GetConfigStruct(path string) (Config, error) {
	var cfg Config

	dataConf, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("couldn't read the file: %w", err)
	}

	err = yaml.Unmarshal(dataConf, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("yaml parsing error: %w", err)
	}

	return cfg, nil
}
