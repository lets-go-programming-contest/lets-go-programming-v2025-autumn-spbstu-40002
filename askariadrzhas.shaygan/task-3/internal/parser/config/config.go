package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input"`
	OutputFile string `yaml:"output"`
}

func ReadConfigPath() *string {
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *configPath == "" {
		panic("no config path provided")
	}

	return configPath
}

func ParseConfig(configPath *string) *Config {
	data, err := os.ReadFile(*configPath)
	if err != nil {
		panic("failed to read config file")
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic("failed to parse config file")
	}

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		panic("invalid config: missing input/output fields")
	}

	return &cfg
}
