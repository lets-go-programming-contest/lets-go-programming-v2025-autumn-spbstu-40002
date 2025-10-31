package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/manyanin.alexander/task-3/internal/config"
	converter "github.com/manyanin.alexander/task-3/internal/currency"
	"github.com/manyanin.alexander/task-3/internal/errors"
	storage "github.com/manyanin.alexander/task-3/internal/json_writer"
	parser "github.com/manyanin.alexander/task-3/internal/xml_parser"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *configPath == "" {
		panic(errors.ErrConfigPathEmpty)
	}

	cfg := config.Load(*configPath)

	if _, err := os.Stat(cfg.InputFile); os.IsNotExist(err) {
		panic(errors.ErrInputFileNotExist.Error() + ": " + cfg.InputFile)
	}

	valCurs := parser.ParseXML(cfg.InputFile)
	currencies := converter.Convert(valCurs)
	storage.SaveToJSON(currencies, cfg.OutputFile)

	fmt.Printf("Successfully processed %d currencies\n", len(currencies))
	fmt.Printf("Data saved to file: %s\n", cfg.OutputFile)
}
