package main

import (
	"fmt"
	"os"

	"github.com/stepanov.alexander/task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config error: %v\n", err)

		return
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
