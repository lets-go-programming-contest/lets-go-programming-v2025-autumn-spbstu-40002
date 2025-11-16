package main

import (
	"flag"

	"github.com/Exam-Play/task-3/internal/config"
	"github.com/Exam-Play/task-3/internal/json"
	"github.com/Exam-Play/task-3/internal/xml"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	data, err := xml.GetCurrencies(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	data.SortOfCurrencies()

	err = json.EncodeJSON(data, cfg.OutputFile)
	if err != nil {
		panic(err)
	}
}
