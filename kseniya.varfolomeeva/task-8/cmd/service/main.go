package main

import (
	"fmt"
	"os"
	"task-8/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Print("error")
		os.Exit(1)
	}
	
	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
