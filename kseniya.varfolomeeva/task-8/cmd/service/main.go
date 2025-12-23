package main

import (
	"fmt"
	"os"
	"task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Print("error")
		os.Exit(1)
	}
	
	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
