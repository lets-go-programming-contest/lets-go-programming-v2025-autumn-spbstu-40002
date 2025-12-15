package main

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	conf "github.com/HuaChenju/task-3/internal/configfile"
	jsonFile "github.com/HuaChenju/task-3/internal/jsonfile"
	xmlsorter "github.com/HuaChenju/task-3/internal/xmlfile/sorter"
	xmlfileparse "github.com/HuaChenju/task-3/internal/xmlfile/xmlparse"
)

var (
	errIncorrectPath     = errors.New("config path is required")
	errIncorrectFormat   = errors.New("unsupported output format")
	errMismatchedTypes   = errors.New("mismatched types")
	errParsingXML        = errors.New("xml parsing failed")
	errWritingOutputFile = errors.New("output file writing failed")
)

const (
	configFlag    string = "config"
	configDefault string = ""
	configHelp    string = "Path to configuration file"
	fomatFlag     string = "output-format"
	formatDefault string = "json"
	formatHelp    string = "Output file format (default: json)"
)

func main() {
	configPath := flag.String(configFlag, configDefault, configHelp)
	outputFormat := flag.String(fomatFlag, formatDefault, formatHelp)

	flag.Parse()

	if *configPath == "" {
		panic(errIncorrectPath)
	}

	if *outputFormat != "json" && *outputFormat != "yaml" && *outputFormat != "xml" {
		panic(errIncorrectFormat)
	}

	cfg, err := conf.GetConfigStruct(*configPath)
	if err != nil {
		panic(err)
	}

	if strings.TrimPrefix(filepath.Ext(cfg.OutputFile), ".") != *outputFormat {
		panic(errMismatchedTypes)
	}

	doc, err := xmlfileparse.GetValCursStruct(cfg.InputFile)
	if err != nil {
		panic(fmt.Errorf("%w: %w", errParsingXML, err))
	}

	xmlsorter.SortValCursByValue(&doc)

	err = jsonFile.WriteToFile(cfg.OutputFile, doc, *outputFormat)
	if err != nil {
		panic(fmt.Errorf("%w: %w", errWritingOutputFile, err))
	}
}
