package main

import (
	"fmt"
	"os"

	"github.com/F0LY/task-8/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
