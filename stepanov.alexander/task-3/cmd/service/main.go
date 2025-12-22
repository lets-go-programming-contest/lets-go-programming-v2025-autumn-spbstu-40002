package main

import (
	"log"

	"stepanov.alexander/task-3/internal/config"
	"stepanov.alexander/task-3/internal/loader"
	"stepanov.alexander/task-3/internal/processor"
	"stepanov.alexander/task-3/internal/writer"
)

func main() {
	cfg := config.ParseFlags()

	xmlData, err := loader.LoadXML(cfg.InputFile)
	if err != nil {
		log.Fatalf("failed to load XML: %v", err)
	}

	rates, err := processor.ProcessXML(xmlData)
	if err != nil {
		log.Fatalf("failed to process XML: %v", err)
	}

	if err := writer.WriteJSON(cfg.OutputFile, rates); err != nil {
		log.Fatalf("failed to write JSON: %v", err)
	}
}
