package main

import (
	"flag"
	"os"

	"github.com/leonid.maryankov/task-3/internal/config"
	"github.com/leonid.maryankov/task-3/internal/parser"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "Path to YAML config file")
	flag.Parse()

	if *cfgPath == "" {
		panic("the path to the configuration file is not specified")
	}

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(cfg.InputFile); err != nil {
		panic(err)
	}

	valutes, err := parser.ParseXML(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	parser.SortValute(valutes)

	if err := parser.SaveToJSON(cfg.OutputFile, valutes); err != nil {
		panic(err)
	}
}
