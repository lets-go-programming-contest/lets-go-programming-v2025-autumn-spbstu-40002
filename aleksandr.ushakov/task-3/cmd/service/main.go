package main

import (
	"github.com/rachguta/task-3/internal/exporter"
	"github.com/rachguta/task-3/internal/parser/config"
	"github.com/rachguta/task-3/internal/parser/xml"
)

func main() {
	// read config path
	configPath := config.ReadConfigPath()

	// parse config file
	cnf := *(config.ParseConfig(configPath))

	// parse xml file
	vc := xml.ParseXML(cnf.InputFile)

	//export to json
	exporter.WriteToJSON(vc.Valute, cnf.OutputFile)
}
