package main

import (
	"flag"

	"github.com/manyanin.alexander/task-3/internal/config"
	converter "github.com/manyanin.alexander/task-3/internal/currency"
	storage "github.com/manyanin.alexander/task-3/internal/json_writer"
	parser "github.com/manyanin.alexander/task-3/internal/xml_parser"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration YAML file")
	flag.Parse()

	cfg := config.Load(*configPath)
	valCurs := parser.ParseXML(cfg.InputFile)
	currencies := converter.Convert(valCurs)
	storage.SaveToJSON(currencies, cfg.OutputFile)
}
