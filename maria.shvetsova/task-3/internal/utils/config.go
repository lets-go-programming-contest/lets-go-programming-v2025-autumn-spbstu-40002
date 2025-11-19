package utils

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	errFileDoesntExist = errors.New("no such file or directory")
	errReadingFile     = errors.New("failed to read the file")
	errOpeningFile     = errors.New("failed to open the file")
	errDecoding        = errors.New("did not find expected key")
	errEmptyFile       = errors.New("file is empty")
	errFileCreating    = errors.New("failed to create a file")
	errInvalidPath     = errors.New("invalid path")
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func parseFlags() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "./config.yaml", "path to config file")

	flag.Parse()

	if configPath == "" {
		return "", errInvalidPath
	}

	return configPath, nil
}

func createOutputFile(path string) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return errFileCreating
	}

	return nil
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

	if _, err = os.Stat(config.OutputFile); os.IsNotExist(err) {
		err = createOutputFile(config.OutputFile)
		if err != nil {
			return nil, fmt.Errorf("creating an output file: %w", err)
		}
	}

	return &config, nil
}
