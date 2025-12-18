package main

import (
	"fmt"
	"log"

	"github.com/lets-go-programming-v2025-autumn-spbstu-40002/victor.kim/task-8/package/config"
)

func main() {
	cfg, err := config.Load(config.ConfigFile)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
