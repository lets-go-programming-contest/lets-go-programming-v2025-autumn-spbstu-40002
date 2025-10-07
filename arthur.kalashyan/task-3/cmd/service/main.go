package main

import (
	"flag"

	"github.com/Expeline/task-3/internal/utils"
)

func main() {
	cfgPath := flag.String("config", "", "")
	flag.Parse()

	if *cfgPath == "" {
		panic("config flag required")
	}

	cfg := utils.LoadConfig(*cfgPath)
	values := utils.ParseCBRXML(cfg.InputFile)
	output := utils.ConvertAndSort(values)
	utils.SaveJSON(cfg.OutputFile, output)
}
