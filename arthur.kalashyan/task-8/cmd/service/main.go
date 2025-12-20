package main

import (
	"fmt"

	"github.com/Expeline/task-8/package/config"
)

func main() {
	cfg := config.Load()
	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}
