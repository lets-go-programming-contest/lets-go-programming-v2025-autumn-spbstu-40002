package main

import (
	"fmt"
	"task-8/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Print("prod error")

		return
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
