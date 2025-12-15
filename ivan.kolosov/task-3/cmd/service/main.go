package main

import (
	"errors"
	"flag"

	cfg "github.com/InsomniaDemon/task-3/internal/config"
	reading "github.com/InsomniaDemon/task-3/internal/readingCurrencies"
	utils "github.com/InsomniaDemon/task-3/internal/utils"
	writing "github.com/InsomniaDemon/task-3/internal/writingCurrencies"
)

var errNoConfigFileProvided = errors.New("no config file was provided")

func main() {
	pathToConfig := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	if *pathToConfig == "" {
		panic(errNoConfigFileProvided)
	}

	config, err := cfg.GetConfig(*pathToConfig)
	if err != nil {
		panic(err)
	}

	curxml, err := reading.GetCurrencies(config.InputFile)
	if err != nil {
		panic(err)
	}

	utils.SortCurrenciesXML(&curxml)

	curjson := utils.GetCurrenciesJSON(curxml)

	err = writing.WriteCurrencies(curjson, config.OutputFile)
	if err != nil {
		panic(err)
	}
}
