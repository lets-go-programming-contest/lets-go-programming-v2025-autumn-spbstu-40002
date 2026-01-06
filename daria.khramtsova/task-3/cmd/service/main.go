package main

import (
	"flag"

	"github.com/hehemka/task-3/internal/utils/config"
	"github.com/hehemka/task-3/internal/utils/json"
	"github.com/hehemka/task-3/internal/utils/xml"
)

func main() {
	configPath := flag.String("config", "config.yaml", "config path")
	flag.Parse()

	conf := config.LoadConfig(*configPath)

	currencies := xml.ReadXML(conf.InputFile)

	json.SortCurrencies(currencies)

	json.WriteJSON(currencies, conf.OutputFile)
}
