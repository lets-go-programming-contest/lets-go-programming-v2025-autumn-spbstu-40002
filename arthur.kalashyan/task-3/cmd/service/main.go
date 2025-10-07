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
	vals := utils.ParseCBRXML(cfg.InputFile)
	outs := utils.ConvertAndSort(vals)
	utils.SaveJSON(cfg.OutputFile, outs)
}
