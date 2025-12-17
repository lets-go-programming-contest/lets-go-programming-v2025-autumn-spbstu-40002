package main

import (
	"fmt"

	"github.com/InsomniaDemon/task-8/internal/config"
)

func main() {
	config, err := config.Load()
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Print(config.Enviroment, " ", config.LogLevel)
}
