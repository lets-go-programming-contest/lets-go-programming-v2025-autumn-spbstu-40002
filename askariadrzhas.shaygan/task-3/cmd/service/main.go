package main

import (
	"log"

	"github.com/XShaygaND/task-3/internal/config"
	"github.com/XShaygaND/task-3/internal/parser"
	"github.com/XShaygaND/task-3/internal/processor"
	"github.com/XShaygaND/task-3/internal/writer"
)

func main() {
	cfgPath := config.ReadConfigPath()
	cfg := config.ParseConfig(cfgPath)

	currencies, err := parser.ParseXML(cfg.InputFile)
	if err != nil {
		log.Fatalf("failed to parse xml: %v", err)
	}

	sorted := processor.OrganizeByRate(currencies)

	if err := writer.WriteJSON(cfg.OutputFile, sorted); err != nil {
		log.Fatalf("failed to write json: %v", err)
	}
}
