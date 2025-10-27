package main

import (
	"errors"
	"flag"
	"log"

	cfg "github.com/ControlShiftEscape/task-3/internal/config"
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
		log.Fatalf("Failed to parse XML file %s: %v", config.Input, err)
	}

}
