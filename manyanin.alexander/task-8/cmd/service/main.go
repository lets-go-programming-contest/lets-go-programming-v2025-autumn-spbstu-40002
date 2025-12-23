package main

import (
	"fmt"
	"os"

	"github.com/manyanin.alexander/task-8/internal/config"
)

func main() {
	appConfig, err := config.LoadAppConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)

		return
	}

	fmt.Print(appConfig.Environment, " ", appConfig.LogLevel)
}
