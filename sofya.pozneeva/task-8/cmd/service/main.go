package main

import (
	"fmt"

	"task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("load config error: %w\n", err)

		return
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
