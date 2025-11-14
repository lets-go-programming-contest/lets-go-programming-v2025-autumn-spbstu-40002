package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/manyanin.alexander/task-3/internal/config"
	"github.com/manyanin.alexander/task-3/internal/currency"
	"github.com/manyanin.alexander/task-3/internal/errors"
	jsonwrite "github.com/manyanin.alexander/task-3/internal/json_writer"
	xmlread "github.com/manyanin.alexander/task-3/internal/xml_parser"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *configPath == "" {
		fmt.Printf("Error: %v\n", errors.ErrConfigPathEmpty)
		return
	}

	cfg := config.Load(*configPath)

	if _, err := os.Stat(cfg.InputFile); os.IsNotExist(err) {
		fmt.Printf("Error: %v: %s\n", errors.ErrInputFileNotExist, cfg.InputFile)
		return
	}

	currencies := xmlread.ReadCurrenciesFromXML(cfg.InputFile)

	if len(currencies) == 0 {
		fmt.Printf("Error: %v\n", errors.ErrNoCurrenciesExtracted)
		return
	}

	sortedCurrencies := currency.SortByValue(currencies)
	jsonwrite.SaveToJSON(sortedCurrencies, cfg.OutputFile)

	fmt.Printf("Successfully processed %d currencies\n", len(sortedCurrencies))
	fmt.Printf("Data saved to file: %s\n", cfg.OutputFile)
}
