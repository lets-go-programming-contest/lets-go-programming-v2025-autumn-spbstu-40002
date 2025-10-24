package main

import (
	"flag"

	"github.com/Exam-Play/task-3/internal/config"
	"github.com/Exam-Play/task-3/internal/jsonfiles"
	"github.com/Exam-Play/task-3/internal/xmlfiles"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	data, err := xmlfiles.GetCurrencies(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	data.SortOfCurrencies()

	err = jsonfiles.EncodeJSON(data, cfg.OutputFile)
	if err != nil {
		panic(err)
	}
}
