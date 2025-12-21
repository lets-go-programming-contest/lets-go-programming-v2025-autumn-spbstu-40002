package main

import (
	"fmt"

	"github.com/Danya-byte/task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("load config error: %w", err)

		return
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
