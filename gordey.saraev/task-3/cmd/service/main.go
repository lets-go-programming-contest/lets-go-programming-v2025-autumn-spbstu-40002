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
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *configPath == "" {
		fmt.Println(errors.ErrConfigPathNotSpecified.Error())

		return
	}

	cfg := config.LoadConfig(*configPath)

	if _, err := os.Stat(cfg.InputFile); os.IsNotExist(err) {
		fmt.Println(errors.ErrInputFileNotExist.Error() + ": " + cfg.InputFile)

		return
	}

	currencies := xmlread.ReadCurrenciesFromXML(cfg.InputFile)

	if len(currencies) == 0 {
		fmt.Println(errors.ErrNoCurrenciesExtracted.Error())

		return
	}

	sortedCurrencies := sort.CurrenciesByValue(currencies)
	jsonwrite.WriteCurrenciesToFile(sortedCurrencies, cfg.OutputFile)

	fmt.Printf("Successfully processed %d currencies\n", len(sortedCurrencies))
	fmt.Printf("Data saved to file: %s\n", cfg.OutputFile)
}
