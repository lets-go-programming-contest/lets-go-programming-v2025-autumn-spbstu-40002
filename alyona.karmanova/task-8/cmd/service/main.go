package main

import (
	"fmt"

	"github.com/HuaChenju/task-8/config"
)

func main() {
	config, err := config.Load()
	if err != nil {
		fmt.Println("Error in config load:", err)

		return
	}

	fmt.Print(config.Environment, " ", config.LogLevel)
}
