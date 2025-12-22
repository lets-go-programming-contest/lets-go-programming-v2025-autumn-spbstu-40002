package main

import (
	"fmt"

	"github.com/t1wt/task-8/package/config"
)

func main() {
	cfg := config.Load()
	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
