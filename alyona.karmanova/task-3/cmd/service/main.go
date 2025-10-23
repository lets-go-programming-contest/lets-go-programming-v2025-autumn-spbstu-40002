package main

import (
	"errors"
	"flag"

	conf "github.com/HuaChenju/task-3/internal/configfile"
	jsonFile "github.com/HuaChenju/task-3/internal/jsonfile"
	xmlfile "github.com/HuaChenju/task-3/internal/xmlfile"
)

var (
	errIncorrectPath = errors.New("config path is required")
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")

	flag.Parse()

	if *configPath == "" {
		panic(errIncorrectPath)
	}

	cfg, err := conf.GetConfigStruct(*configPath)
	if err != nil {
		panic(err)
	}

	doc, err := xmlfile.GetValCursStruct(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	xmlfile.SortValCursByValue(&doc)

	err = jsonFile.WriteJSONToFile(cfg.OutputFile, doc)
	if err != nil {
		panic(err)
	}
}
