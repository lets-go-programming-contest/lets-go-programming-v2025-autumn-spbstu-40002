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
	cfgPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *cfgPath == "" {
		panic("missing --config flag")
	}

	return cfgPath
}

func ParseConfig(cfgPath *string) *Config {
	data, err := os.ReadFile(*cfgPath)
	if err != nil {
		panic("failed to read config file")
	}

	var cnf Config
	if err := yaml.Unmarshal(data, &cnf); err != nil {
		panic("failed to parse config file")
	}

	if cnf.InputFile == "" || cnf.OutputFile == "" {
		panic("invalid config: missing input/output fields")
	}

	return &cnf
}
