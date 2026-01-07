package main

import (
	"fmt"
	"log"

	"github.com/ControlShiftEscape/task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.Environment, cfg.LogLevel)
}
