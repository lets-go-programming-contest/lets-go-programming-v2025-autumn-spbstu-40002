package main

import (
	"errors"
	"flag"
)

var (
	errNoConfigFileProvided = errors.New("no config file was provided")
)

func main() {
	pathToConfig := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *pathToConfig == "" {
		panic(errNoConfigFileProvided)
	}

	config, err := cfg.getConfig(*pathToConfig)
	if err != nil {
		panic(err)
	}

	curxml, err := reading.getCurrencies(config.inputFile)
	if err != nil {
		panic(err)
	}

	utils.sortCurrenciesXML(curxml)

	curjson := utils.getCurrenciesJSON(curxml)

	err = writing.writeCurrencies(curjson, config.outputFile)
	if err != nil {
		panic(err)
	}
}
