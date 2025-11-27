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
	configPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	if configPath == nil || *configPath == "" {
		panic(ErrEmptyConfigFlag)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	doc, err := cbr.LoadXML(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	valutes, err := convert.Sort(doc)
	if err != nil {
		panic(err)
	}

	if err := output.WriteJSON(cfg.OutputFile, valutes); err != nil {
		panic(err)
	}
}
