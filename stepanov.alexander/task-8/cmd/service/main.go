package main

import (
	"fmt"
	"os"

	"stepanov.alexander/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config error: %v", err)

		return
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}