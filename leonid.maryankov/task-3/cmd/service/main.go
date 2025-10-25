package main

import (
	"flag"
	"log"
	"os"

	"github.com/leonid.maryankov/task-3/internal/config"
	"github.com/leonid.maryankov/task-3/internal/parser"
)

func main() {
	cfgPath := flag.String("config", "", "Path to YAML config file")
	flag.Parse()

	if *cfgPath == "" {
		log.Fatal("the path to the configuration file is not specified")
	}

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(cfg.InputFile); err != nil {
		log.Fatal(err)
	}

	valutes, err := parser.ParseXML(cfg.InputFile)
	if err != nil {
		log.Fatal(err)
	}

	parser.SortValute(valutes)

	if err := parser.SaveToJson(cfg.OutputFile, valutes); err != nil {
		log.Fatal(err)
	}
}
