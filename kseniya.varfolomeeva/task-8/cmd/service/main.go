package main

import (
	"fmt"
	"task-8/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Err with the loading file: ", err)
		return
	}
	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
