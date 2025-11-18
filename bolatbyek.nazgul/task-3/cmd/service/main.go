package main

import (
	"flag"
	"fmt"

	"github.com/Nazkaaa/task-3/internal/config"
	"github.com/Nazkaaa/task-3/internal/jsonwriter"
	"github.com/Nazkaaa/task-3/internal/models"
	"github.com/Nazkaaa/task-3/internal/xmlparser"
)

func run(configPath string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	err = config.EnsureOutputDir(cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	valCurs, err := xmlparser.ParseXML(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("error parsing XML: %w", err)
	}

	currencies := models.ConvertAndSort(valCurs)

	outputCurrencies := make([]interface{}, 0, len(currencies))
	for _, c := range currencies {
		outputCurrencies = append(outputCurrencies, c)
	}

	err = jsonwriter.SaveToJSON(outputCurrencies, cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("error saving to JSON: %w", err)
	}

	fmt.Printf("Successfully processed %d currencies and saved to %s\n", len(currencies), cfg.OutputFile)

	return nil
}

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the YAML configuration file")
	flag.Parse()

	if err := run(*configPath); err != nil {
		panic(err.Error())
	}
}
