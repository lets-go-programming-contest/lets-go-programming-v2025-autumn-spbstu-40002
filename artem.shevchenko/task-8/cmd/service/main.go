package main

import (
	"fmt"

	"github.com/slendycs/go-lab-8/internal/config"
)

func main() {
	config, err := config.Load()
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(config.Environment, " ", config.LogLevel)
}
