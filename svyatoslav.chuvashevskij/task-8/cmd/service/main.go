package main

import (
	"fmt"
	"log"

	"github.com/Svyatoslav2324/task-8/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
