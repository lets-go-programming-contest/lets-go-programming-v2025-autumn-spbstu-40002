package main

import (
	"flag"

	"github.com/xanaoxx/task-3/internal/config"
	"github.com/xanaoxx/task-3/internal/json"
	"github.com/xanaoxx/task-3/internal/xml"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	appConfig, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	currencyRecords, err := xml.LoadCurrencies(appConfig.InputFile)
	if err != nil {
		panic(err)
	}

	currencyRecords.SortByValue()

	err = json.SaveCurrencies(currencyRecords, appConfig.OutputFile)
	if err != nil {
		panic(err)
	}
}
