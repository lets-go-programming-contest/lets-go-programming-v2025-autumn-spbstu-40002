package main

import (
	"fmt"

	"github.com/Exam-Play/task-8/internal/config"
)

func main() {
	cfg, err := config.LoadCfg()
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(cfg.Environment, cfg.LogLevel)
}
