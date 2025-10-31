package main

import (
	"flag"

	"github.com/t1wt/task-3/internal/utils"
)

func main() {
	cfg := flag.String("config", "config.yaml", "")
	flag.Parse()
	if err := utils.Execute(*cfg); err != nil {
		panic(err)
	}
}
