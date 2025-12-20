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
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)

		return
	}

	if _, err := os.Stat(cfg.InputFile); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, errors.ErrInputFileNotExist.Error()+": "+cfg.InputFile)

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
