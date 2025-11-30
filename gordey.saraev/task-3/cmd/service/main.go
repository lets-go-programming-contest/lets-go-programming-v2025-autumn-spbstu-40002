package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/F0LY/task-3/internal/config"
	"github.com/F0LY/task-3/internal/errors"
	"github.com/F0LY/task-3/internal/jsonwrite"
	"github.com/F0LY/task-3/internal/sort"
	"github.com/F0LY/task-3/internal/xmlread"
)

func main() {
	// add default.
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, errors.ErrConfigPathNotSpecified.Error()+": "+*configPath)

		return
	}

	cfg := config.LoadConfig(*configPath)

	if _, err := os.Stat(cfg.InputFile); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, errors.ErrInputFileNotExist.Error()+": "+cfg.InputFile)

		return
	}

	currencies := xmlread.ReadCurrenciesFromXML(cfg.InputFile)
	if len(currencies) == 0 {
		fmt.Fprintln(os.Stderr, errors.ErrNoCurrenciesExtracted.Error())

		return
	}

	sortedCurrencies := sort.CurrenciesByValue(currencies)
	jsonwrite.WriteCurrenciesToFile(sortedCurrencies, cfg.OutputFile)

	fmt.Printf("Successfully processed %d currencies\n", len(sortedCurrencies))
	fmt.Printf("Data saved to file: %s\n", cfg.OutputFile)
}
