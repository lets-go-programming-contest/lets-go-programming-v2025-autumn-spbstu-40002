package main

import (
	"log"

	"stepanov.alexander/task-3/internal/config"
	"stepanov.alexander/task-3/internal/loader"
	"stepanov.alexander/task-3/internal/processor"
	"stepanov.alexander/task-3/internal/writer"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		log.Fatalf("input-file or output-file is empty")
	}

	valCurs, err := loader.LoadXML(cfg.InputFile)
	if err != nil {
		log.Fatalf("failed to load XML: %v", err)
	}

	rates, err := processor.ProcessXML(valCurs)
	if err != nil {
		log.Fatalf("failed to process XML: %v", err)
	}

	if err := writer.WriteJSON(cfg.OutputFile, rates); err != nil {
		log.Fatalf("failed to write JSON: %v", err)
	}
}
