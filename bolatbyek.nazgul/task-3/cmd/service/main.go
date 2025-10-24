package main

import (
	"flag"
	"fmt"
	"github.com/bolatbyek/task-3/internal/config"
	"github.com/bolatbyek/task-3/internal/currency"
)

func main() {
	// Parse command line arguments
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *configPath == "" {
		panic("Config file path is required")
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Process currency data
	processor := currency.NewProcessor()
	
	// Parse XML
	valCurs, err := processor.ParseXML(cfg.InputFile)
	if err != nil {
		panic("Failed to parse XML: " + err.Error())
	}

	// Sort currencies by value (descending)
	sortedCurrencies := processor.SortByValue(valCurs.Currencies)

	// Save to JSON
	err = processor.SaveToJSON(sortedCurrencies, cfg.OutputFile)
	if err != nil {
		panic("Failed to save JSON: " + err.Error())
	}

	fmt.Printf("Successfully processed %d currencies and saved to %s\n", 
		len(sortedCurrencies), cfg.OutputFile)
}
