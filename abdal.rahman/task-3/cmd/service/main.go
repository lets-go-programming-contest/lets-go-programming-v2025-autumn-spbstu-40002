package main

import (
	"flag"

	"github.com/xkoex/task-3/internal/config"
	"github.com/xkoex/task-3/internal/jsonutils"
	"github.com/xkoex/task-3/internal/xmlutils"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	conf := config.LoadConfig(*configPath)

	currencies := xmlutils.ReadXML(conf.InputFile)

	jsonutils.SortCurrencies(currencies)

	jsonutils.WriteJSON(currencies, conf.OutputFile)
}
