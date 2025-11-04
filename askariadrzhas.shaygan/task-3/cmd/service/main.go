package main

import (
	"log"

	"github.com/XShaygaND/task-3/internal/config"
	"github.com/XShaygaND/task-3/internal/parser"
	"github.com/XShaygaND/task-3/internal/processor"
	"github.com/XShaygaND/task-3/internal/writer"
)

func main() {
	settings, err := config.LoadSettings()
	if err != nil {
		log.Fatalf("failed to load settings: %v", err)
	}

	data, err := parser.ParseXML(settings.InputPath)
	if err != nil {
		log.Fatalf("failed to parse xml: %v", err)
	}

	sorted := processor.OrganizeByRate(data)

	if err := writer.WriteJSON(settings.OutputPath, sorted); err != nil {
		log.Fatalf("failed to write json: %v", err)
	}
}
