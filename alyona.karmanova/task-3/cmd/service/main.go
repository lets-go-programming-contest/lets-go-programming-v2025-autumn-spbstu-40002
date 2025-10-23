package main

import (
	"errors"
	"flag"
	"path/filepath"
	"strings"

	conf "task-3/internal/configfile"
	outputFile "task-3/internal/outputfile"
	xmlfile "task-3/internal/xmlfile"
)

var (
	errIncorrectPath   = errors.New("config path is required")
	errIncorrectFormat = errors.New("unsupported output format")
	errMismatchedTypes = errors.New("mismatched types")
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	outputFormat := flag.String("output-format", "json", "Output file format (default: json)")

	flag.Parse()

	if *configPath == "" {
		panic(errIncorrectPath)
	}

	if *outputFormat != "" && *outputFormat != "json" && *outputFormat != "yaml" && *outputFormat != "xml" {
		panic(errIncorrectFormat)
	}

	cfg, err := conf.GetConfigStruct(*configPath)
	if err != nil {
		panic(err)
	}

	if strings.TrimPrefix(filepath.Ext(cfg.OutputFile), ".") != *outputFormat {
		panic(errMismatchedTypes)
	}

	doc, err := xmlfile.GetValCursStruct(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	xmlfile.SortValCursByValue(&doc)

	err = outputFile.WriteToFile(cfg.OutputFile, doc, *outputFormat)
	if err != nil {
		panic(err)
	}
}
