package main

import (
	"errors"
	"flag"

	cfg "github.com/ControlShiftEscape/task-3/internal/config"
)

var errInvalidConfigFile = errors.New("invalid config file")

func main() {

	configPath := flag.String("config", "", "Config directory")
	flag.Parse()

	if *configPath == "" {
		panic(errInvalidConfigFile)
	}

	config, err := cfg.ParseConfig(*configPath)
	if err != nil {
		panic(err)
	}

}
