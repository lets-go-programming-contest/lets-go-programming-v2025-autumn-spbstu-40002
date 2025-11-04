package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppSettings struct {
	SourcePath string `yaml:"inputFile"`
	TargetPath string `yaml:"outputFile"`
}

func LoadSettings() *AppSettings {
	configPath := flag.String("config", "", "Configuration file path")
	flag.Parse()

	if *configPath == "" {
		panic("configuration path not provided")
	}

	fileData, err := os.ReadFile(*configPath)
	if err != nil {
		panic(fmt.Sprintf("cannot read config file: %v", err))
	}

	var settings AppSettings
	if err := yaml.Unmarshal(fileData, &settings); err != nil {
		panic(fmt.Sprintf("invalid config format: %v", err))
	}

	if settings.SourcePath == "" || settings.TargetPath == "" {
		if _, err := os.Stat(settings.SourcePath); os.IsNotExist(err) {
			panic(fmt.Sprintf("cannot read source file: %v", err))
		}

		panic("missing required paths in configuration")
	}

	return &settings
}
