package main

import (
	"github.com/Tsahaev/task-3/internal/exporter"
	"github.com/Tsahaev/task-3/internal/parser/config"
	"github.com/Tsahaev/task-3/internal/parser/xmlparser"
)

func main() {
	confPath := config.ReadConfigPath()

	conf := config.ParseConfig(confPath)

	xmlData := xmlparser.ParseXML(conf.InputFile)

	exporter.WriteToJSON(xmlData.Valute, conf.OutputFile)
}
