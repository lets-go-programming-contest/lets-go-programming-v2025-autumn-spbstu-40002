package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type AppSettings struct {
	SourcePath string `yaml:"input-file"`
	TargetPath string `yaml:"output-file"`
}

func LoadSettings() *AppSettings {
	configPath := flag.String("config", "", "Configuration file path")
	flag.Parse()

	if *configPath == "" {
		panic("configuration path not provided")
	}

	fileData, err := os.ReadFile(*configPath)

	if err != nil {
		panic("cannot read config file: " + err.Error())
	}

	var settings AppSettings
	err = yaml.Unmarshal(fileData, &settings)

	if err != nil {
		panic("invalid config format: " + err.Error())
	}

	if settings.SourcePath == "" || settings.TargetPath == "" {
		panic("missing required paths in configuration")
	}

	return &settings
}
