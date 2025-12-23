package main

import (
	"fmt"

	"platon.karyakin/task-8/internal/config"
)

func main() {
	cnfg, err := config.Initial()
	if err != nil {
		fmt.Println("failed config: %w", err)

		return
	}

	cnfg.PrintToConfig()
}
