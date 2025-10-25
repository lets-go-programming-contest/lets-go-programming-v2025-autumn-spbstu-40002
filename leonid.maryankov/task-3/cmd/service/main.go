package main

import (
	"flag"
	"log"
	"os"

	"github.com/leonid.maryankov/task-3/internal/config"
	"github.com/leonid.maryankov/task-3/internal/parser"
)

func main() {
	configPath := flag.String("config", "", "Path")
	flag.Parse()

	if *configPath == "" {
		panic("The path to the configuration file is not specified")
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(cfg.InputFile); os.IsNotExist(err) {
		log.Fatalln("The input XML file was not found" + cfg.InputFile)
	}

	valutes, err := parser.ParseXML(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	parser.SortValute(valutes)

	if err := parser.SaveToJson(cfg.OutputFile, valutes); err != nil {
		panic(err)
	}

	log.Println("Saved in" + cfg.OutputFile)
}
