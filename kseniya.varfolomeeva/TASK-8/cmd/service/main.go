package main

import (
	"fmt"
	"task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}
	
	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}
