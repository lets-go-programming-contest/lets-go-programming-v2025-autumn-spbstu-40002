package main

import (
	"fmt"
	"os"

	"github.com/F0LY/task-8/internal/config"
)

func main() {
	cfg, err := config.LoadAppConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}
