package main

import (
	"flag"

	"github.com/slendycs/go-lab-3/internal/config"
	merr "github.com/slendycs/go-lab-3/internal/myerrors"
	"github.com/slendycs/go-lab-3/internal/parsers"
	"github.com/slendycs/go-lab-3/internal/utils"
)

func main() {
	// Read config path from cli flags.
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		panic(merr.ErrNoConfigFileProvided)
	}

	// Read config from YAML.
	var cfg config.Config
	config.ReadConfig(configPath, &cfg)

	// Read Valute data from XML.
	valData := new(parsers.ValStruct)
	parsers.ReadXML(cfg.InputFile, valData)

	// Sort Valute data.
	utils.SortVal(valData)

	// Write JSON file of Valutes
	parsers.WriteJSON(cfg.OutputFile, valData)
}
