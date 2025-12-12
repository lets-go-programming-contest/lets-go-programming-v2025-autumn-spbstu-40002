package main

import (
	"config-example/config"
	"log"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	cfg.PrintConfig()
}
