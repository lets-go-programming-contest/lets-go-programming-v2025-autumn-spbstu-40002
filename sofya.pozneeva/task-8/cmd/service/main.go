package main

import (
	"fmt"

	"task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("load config error: ", err)

		return
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
