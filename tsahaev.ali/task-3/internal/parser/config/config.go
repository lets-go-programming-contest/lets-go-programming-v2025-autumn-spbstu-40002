package config

import (
	"flag"
	"os"

	"github.com/Tsahaev/task-3/internal/myerrors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func ReadConfigPath() *string {
	confPath := flag.String("config", "", "path to config YAML file")
	flag.Parse()

	if *confPath == "" {
		panic(myerrors.ErrConfigPath)
	}
	return confPath
}

func ParseConfig(path *string) *Config {
	data, err := os.ReadFile(*path)
	if err != nil {
		panic(myerrors.ErrConfigRead)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil || cfg.InputFile == "" || cfg.OutputFile == "" {
		panic(myerrors.ErrConfigParse)
	}

	return &cfg
}
