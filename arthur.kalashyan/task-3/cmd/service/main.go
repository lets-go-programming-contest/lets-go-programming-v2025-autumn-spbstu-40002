package main

import (
	"flag"

	utils "github.com/Expeline/task-3/internal/utils"
)

func main() {
	configPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		panic("config path not provided")
	}

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
