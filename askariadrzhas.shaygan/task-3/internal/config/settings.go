package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type AppSettings struct {
	SourcePath string `yaml:"inputFile"`
	TargetPath string `yaml:"outputFile"`
}

func LoadSettings() *AppSettings {
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *configPath == "" {
		panic("configuration path not provided")
	}

	fileData, err := os.ReadFile(*configPath)
	if err != nil {
		panic("cannot read config file: " + err.Error())
	}

	var settings AppSettings
	if err = yaml.Unmarshal(fileData, &settings); err != nil {
		panic("invalid config format: " + err.Error())
	}

	if settings.SourcePath == "" || settings.TargetPath == "" {
		panic("missing required paths in configuration")
	}

	return &settings
}
