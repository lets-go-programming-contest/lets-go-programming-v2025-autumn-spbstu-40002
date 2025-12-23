package main

import (
	"fmt"
	"task-8/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		if cfg != nil && cfg.Environment != "" {
			fmt.Print(cfg.Environment, " error")
		} else {
			fmt.Print("error")
		}
		return
	}
	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
