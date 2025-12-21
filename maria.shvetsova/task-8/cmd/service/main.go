package main

import (
	"fmt"

	"github.com/ummmsh/task-8/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
