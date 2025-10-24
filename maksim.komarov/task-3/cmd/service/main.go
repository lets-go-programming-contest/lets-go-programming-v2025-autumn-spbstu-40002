package main

import (
	"errors"
	"flag"

	"github.com/megurumacabre/task-3/internal/cbr"
	"github.com/megurumacabre/task-3/internal/config"
	"github.com/megurumacabre/task-3/internal/convert"
	"github.com/megurumacabre/task-3/internal/output"
)

var ErrEmptyConfigFlag = errors.New("config flag is empty")

func main() {
	configPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		panic(ErrEmptyConfigFlag)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	doc, err := cbr.ReadFile(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	out, err := convert.MapAndSort(doc)
	if err != nil {
		panic(err)
	}
	
	if err := output.WriteJSON(cfg.OutputFile, out); err != nil {
		panic(err)
	}
}
