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

	fmt.Printf("DEBUG: Config path: %s\n", *configPath)

	appConfig, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("DEBUG: Config error: %v\n", err)
		panic(err)
	}

	fmt.Printf("DEBUG: Input file: %s\n", appConfig.InputFile)

	currencyRecords, err := xml.LoadCurrencies(appConfig.InputFile)
	if err != nil {
		fmt.Printf("DEBUG: XML error: %v\n", err)
		panic(err)
	}

	currencyRecords.SortByValue()

	err = json.SaveCurrencies(currencyRecords, appConfig.OutputFile)
	if err != nil {
		fmt.Printf("DEBUG: JSON error: %v\n", err)
		panic(err)
	}
}
