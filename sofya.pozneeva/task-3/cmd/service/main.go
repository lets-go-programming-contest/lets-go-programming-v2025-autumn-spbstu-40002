package main

import (
	"task-3/internal/utils"
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

	valutesForJSON, err := utils.GetValutesForJSON((*utils.Valutes)(&valCurs.Valutes))
	if err != nil {
		panic(err)
	}

	utils.JSONWrite(config, valutesForJSON)
}
