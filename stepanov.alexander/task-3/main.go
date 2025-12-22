package main

import (
	"flag"
	"fmt"
	"os"

	"stepanov.alexander/task-3/internal/config"
	"stepanov.alexander/task-3/internal/currency"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to YAML configuration file")
	flag.Parse()

	if configPath == "" {
		fmt.Fprintf(os.Stderr, "Error: Flag --config must be set\n")
		os.Exit(1)
	}

	cfg, err := config.ReadSettings(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to read config: %v\n", err)
		os.Exit(1)
	}

	exchangeRates, err := currency.FetchCurrencyRates(cfg.InputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to fetch currency rates: %v\n", err)
		os.Exit(1)
	}

	if err := currency.OrderByExchange(exchangeRates); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to sort currencies: %v\n", err)
		os.Exit(1)
	}

	if err := currency.ExportToJSON(cfg.OutputFile, exchangeRates); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to export JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %d currencies. Output: %s\n",
		len(exchangeRates.Currencies), cfg.OutputFile)
}
