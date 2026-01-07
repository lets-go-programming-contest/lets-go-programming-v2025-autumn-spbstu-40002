package main

import (
	"flag"

	"github.com/braab/lets-go-programming-v2025-autumn-spbstu-40002/abdal.rahman/task-3/internal/config"
	"github.com/braab/lets-go-programming-v2025-autumn-spbstu-40002/abdal.rahman/task-3/internal/parser"
	"github.com/braab/lets-go-programming-v2025-autumn-spbstu-40002/abdal.rahman/task-3/internal/parser/utils"
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
