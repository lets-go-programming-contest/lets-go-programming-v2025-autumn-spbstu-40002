package main

import (
	"flag"
	"fmt"

	"github.com/platon.karyakin/task-3/internal/bank"
	"github.com/platon.karyakin/task-3/internal/config"
	"github.com/platon.karyakin/task-3/pkg/must"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	appConfig, errLoadConfig := config.LoadConfigFromFile(*configPath)
	must.Must("load config file", errLoadConfig)

	exchangeRates, errLoadXML := bank.LoadFromXML(appConfig.Input)
	must.Must("load exchange XML", errLoadXML)

	if errSave := exchangeRates.WriteJSONFile(appConfig.Output); errSave != nil {
		must.Must("save output json", errSave)
	}

	fmt.Println("ok")
}
