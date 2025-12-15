package main

import (
	"flag"

	"github.com/Expeline/task-3/internal/utils"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := utils.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	data, err := utils.ReadXML(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	sorted := utils.SortCurrencies(data.Currencies)

	if err := utils.SaveToJSON(sorted, cfg.OutputFile); err != nil {
		panic(err)
	}
}
