package main

import (
	"flag"

	"danila.onitshuk/task-3/internal/config"
	"danila.onitshuk/task-3/internal/parser"
	"danila.onitshuk/task-3/internal/parser/utils"
)

func main() {
	var (
		configPath string
		cfg        config.Config
	)

	flag.StringVar(&configPath, "config", "", "path to config file")

	flag.Parse()

	if configPath == "" {
		panic(config.ErrPath)
	}

	config.NewConfig(configPath, &cfg)

	inputData := new(parser.XMLData)
	parser.ReadXML(cfg.InputFile, inputData)

	outputData := utils.ToCurrency(inputData.Valute)
	utils.SortVal(outputData)

	parser.WriteJSON(cfg.OutputFile, outputData)
}
