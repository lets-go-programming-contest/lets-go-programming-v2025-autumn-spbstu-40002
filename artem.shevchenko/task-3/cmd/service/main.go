package main

import (
	"flag"

	"github.com/slendycs/go-lab-3/internal/config"
	merr "github.com/slendycs/go-lab-3/internal/myerrors"
	"github.com/slendycs/go-lab-3/internal/json"
	"github.com/slendycs/go-lab-3/internal/xml"
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
	valData := new(xml.ValCurs)
	xml.ReadXML(cfg.InputFile, valData)

	// Sort Valute data.
	xml.SortVal(valData)

	// Write JSON file of Valutes
	json.WriteJSONData(cfg.OutputFile, valData)
}
