package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func ParseFlags() (*Config, error) {
	configPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()

	cfg := &Config{
		InputFile:  "",
		OutputFile: "",
	}

	if configPath == nil || *configPath == "" {
		// допускаем запуск без конфига: вернём пустой cfg и nil
		return cfg, nil
	}

	data, err := os.ReadFile(*configPath)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
