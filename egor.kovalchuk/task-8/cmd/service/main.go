package main

import (
	"fmt"

	"egor.kovalchuk/task-8/config"
)

func main() {
	config, err := config.ParseConfig()
	if err != nil {
		fmt.Println("Error: ", err)

		return
	}

	fmt.Print(config.Environment, " ", config.LogLevel)
}
