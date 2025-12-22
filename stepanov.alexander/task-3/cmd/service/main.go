package main

import (
	"log"

	"github.com/stepanov.alexander/task-3/pkg/config"
	"github.com/stepanov.alexander/task-3/pkg/loader"
	"github.com/stepanov.alexander/task-3/pkg/processor"
	"github.com/stepanov.alexander/task-3/pkg/writer"
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

	err = writer.WriteJSON(cfg.OutputFile, rates)
	if err != nil {
		log.Fatalf("failed to write JSON: %v", err)
	}
}
