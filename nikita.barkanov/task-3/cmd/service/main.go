package main

import (
	"errors"
	"flag"

	cfg "github.com/ControlShiftEscape/task-3/internal/config"
	jsonw "github.com/ControlShiftEscape/task-3/internal/jsonwriter"
	xmlp "github.com/ControlShiftEscape/task-3/internal/xmlparser"
)

var errInvalidConfigFile = errors.New("invalid config file")

func main() {

	configPath := flag.String("config", "config.yaml", "Config directory")
	flag.Parse()

	if *configPath == "" {
		panic(errInvalidConfigFile)
	}

	config, err := cfg.ParseConfig(*configPath)
	if err != nil {
		panic(err)
	}

	curs, err := xmlp.ParseXML(config.Input)
	if err != nil {
		//log.Printf("Input file:  %s", config.Input)
	}

	if err := jsonw.WriteSortedReducedJSON(curs, config.Output); err != nil {
		//log.Fatalf("Failed to write JSON to %s: %v", config.Output, err)
	}

}
