package main

import (
	"flag"

	"github.com/slendycs/go-lab-3/internal/config"
	"github.com/slendycs/go-lab-3/internal/parsers"
	"github.com/slendycs/go-lab-3/internal/utils"
)

func main() {
	var (
		cfg        config.Config
		configPath string
	)

	// Read config path from cli flags.
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		configPath = "config.yaml" // Default config path
	}

	// Read config from YAML.
	err := config.ReadConfig(configPath, &cfg)
	if err != nil {
		panic(err)
	}

	// Read Valute data from XML.
	valData := new(parsers.ValStruct)
	parsers.ReadXML(cfg.InputFile, valData)

	// Sort Valute data.
	utils.SortVal(valData)

	// Write JSON file of Valutes
	parsers.WriteJSON(cfg.OutputFile, valData)
}
