package config

import (
	"flag"
)

type Config struct {
	InputFile  string
	OutputFile string
}

func ParseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.InputFile, "config", "", "path to input XML file")
	flag.StringVar(&cfg.OutputFile, "output-file", "", "path to output JSON file")
	flag.Parse()

	return cfg
}
