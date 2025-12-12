package main

import (
	"fmt"
	"os"

	"github.com/F0LY/task-8/internal/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "when loading config: %v\n", err)

		return
	}

	fmt.Print(config.Environment, " ", config.LogLevel)
}
