package main

import (
	"errors"
	"fmt"

	"github.com/HuaChenju/task-8/config"
)

var ErrUnmarshalConfig = errors.New("failed to config load")

func main() {
	config, err := config.Load()
	if err != nil {
		fmt.Println("Error in config load:", err)

		return
	}

	fmt.Print(config.Environment, " ", config.LogLevel)
}
