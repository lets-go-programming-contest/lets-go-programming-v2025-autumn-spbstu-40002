package main

import (
	"log"

	"github.com/XShaygaND/task-3/internal/parser"
	"github.com/XShaygaND/task-3/internal/parser/config"
	"github.com/XShaygaND/task-3/internal/processor"
	"github.com/XShaygaND/task-3/internal/writer"
)

func main() {
	configPath := config.ReadConfigPath()
	cfg := config.ParseConfig(configPath)

	items, err := parser.ParseXML(cfg.InputFile)
	if err != nil {
		log.Fatalf("failed to parse XML: %v", err)
	}

	sortedItems := processor.SortData(items)

	if err := writer.WriteJSON(cfg.OutputFile, sortedItems); err != nil {
		log.Fatalf("failed to write JSON: %v", err)
	}
}
