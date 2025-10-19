package main

import (
	"flag"
	"fmt"

	"github.com/victor.kim/task-3/internal/bank"
	"github.com/victor.kim/task-3/internal/config"
	"github.com/victor.kim/task-3/pkg/must"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.ParseFile(*configPath)
	must.Must("parse config", err)

	bankData, err := bank.ParseFileXML(cfg.Input)
	must.Must("parse input-file", err)

	if err := bankData.EncodeFileJSON(cfg.Output); err != nil {
		must.Must("encode bank", err)
	}

	fmt.Println("ok")
}
