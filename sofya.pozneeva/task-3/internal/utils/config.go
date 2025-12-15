package utils

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	errInvalidPath     = errors.New("invalid yaml path")
	errFileDoesntExist = errors.New("no such file or directory")
	errReadingFile     = errors.New("failed to read the file")
	errOpeningFile     = errors.New("failed to open the file")
	errDecoding        = errors.New("did not find expected key")
	errEmptyFile       = errors.New("file is empty")
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func parseFlags() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to config file")

	flag.Parse()

	if configPath == "" {
		return "", errInvalidPath
	}

	return configPath, nil
}

func GetConfig() (*Config, error) {
	path, err := parseFlags()
	if err != nil {
		return nil, fmt.Errorf("parse flags: %w", err)
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		return nil, errFileDoesntExist
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errReadingFile
	}

	if len(data) == 0 {
		return nil, errEmptyFile
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, errDecoding
	}

	return &config, nil
}
