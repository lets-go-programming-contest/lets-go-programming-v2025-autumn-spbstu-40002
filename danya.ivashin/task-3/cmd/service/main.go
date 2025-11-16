package main

import (
	"flag"

	"github.com/Danya-byte/task-3/internal/bank"
	"github.com/Danya-byte/task-3/internal/config"
	"github.com/Danya-byte/task-3/pkg/must"
)

func main() {
	configPath := flag.String("config", "config.yaml", "select path to config file")
	flag.Parse()

	config, err := config.ParseFile(*configPath)

	must.Must("parse config", err)

	bank, err := bank.ParseFileXML(config.Input)

	must.Must("parse input-file", err)
	must.Must("encode bank", bank.EncodeFileJSON(config.Output))
}
