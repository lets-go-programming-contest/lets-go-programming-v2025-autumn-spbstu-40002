package config

import (
	"flag"
	"os"

	"github.com/rachguta/task-3/internal/myerrors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func ReadConfigPath() *string {
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *configPath == "" {
		panic(myerrors.ErrConfigPath)
	}

	return configPath
}

func ParseConfig(configPath *string) *Config {
	// read config
	data, err := os.ReadFile(*configPath)
	if err != nil {
		panic(myerrors.ErrConfigRead)
	}
	// parse config
	var cnf Config

	err = yaml.Unmarshal(data, &cnf)
	if err != nil || cnf.InputFile == "" || cnf.OutputFile == "" {
		panic(myerrors.ErrConfigParse)
	}

	return &cnf
}
