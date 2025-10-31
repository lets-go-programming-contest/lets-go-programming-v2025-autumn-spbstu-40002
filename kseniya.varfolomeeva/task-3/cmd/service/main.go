package main

import (
	"flag"
	"fmt"

	"github.com/xanaoxx/task-3/internal/config"
	"github.com/xanaoxx/task-3/internal/json"
	"github.com/xanaoxx/task-3/internal/xml"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	appConfig, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("DEBUG: Input file: %s\n", appConfig.InputFile)

	currencyRecords, err := xml.LoadCurrencies(appConfig.InputFile)
	if err != nil {
		panic(err)
	}

	currencyRecords.SortByValue()

	err = json.SavCurrencies(currencyRecords, appConfig.OutputFile)
	if err != nil {
		panic(err)
	}
}
