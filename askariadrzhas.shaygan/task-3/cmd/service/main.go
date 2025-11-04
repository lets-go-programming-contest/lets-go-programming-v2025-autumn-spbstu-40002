package main

import (
	"log"

	"github.com/XShaygaND/task-3/internal/parser"
	"github.com/XShaygaND/task-3/internal/parser/config"
	"github.com/XShaygaND/task-3/internal/processor"
	"github.com/XShaygaND/task-3/internal/writer"
)

func main() {
	cfgPath := config.ReadConfigPath()
	cfg := config.ParseConfig(cfgPath)

	data, err := parser.ParseXML(cfg.InputFile)
	if err != nil {
		log.Fatalf("failed to parse XML: %v", err)
	}

	sorted := processor.OrganizeByRate(data)

	err = writer.WriteJSON(cfg.OutputFile, sorted)
	if err != nil {
		log.Fatalf("failed to write JSON: %v", err)
	}
}
