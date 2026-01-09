package main

import (
	"fmt"

	"github.com/xkoex/task-3/internal/core"
)

func main() {
	cfg, err := core.LoadConfig()
	if err != nil {
		panic(err)
	}

	data, err := core.GetCurrencyData(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	data.SortByValue()

	jsonBytes, err := core.PrepareJSON((*core.Valutes)(&data.Valutes))
	if err != nil {
		panic(err)
	}

	err = core.JSONWrite(cfg, jsonBytes)
	if err != nil {
		panic(err)
	}

	fmt.Println("JSON file created successfully:", cfg.OutputFile)
}
