package main

import (
	"github.com/ummmsh/task-3/internal/utils"
)

func main() {
	config, err := utils.GetConfig()
	if err != nil {
		panic(err)
	}

	valCurs, err := utils.GetCurrencyData(config.InputFile)
	if err != nil {
		panic(err)
	}

	valCurs.SortByValue()

	outputData, err := valCurs.ConvertToOutput()
	if err != nil {
		panic(err)
	}

	err = utils.WriteToJSON(outputData, config.OutputFile)
	if err != nil {
		panic(err)
	}
}
