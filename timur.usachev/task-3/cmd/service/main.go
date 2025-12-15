package main

import (
	"flag"

	"github.com/t1wt/task-3/internal/utils"
)

func main() {
	cfg := flag.String("config", "config.yaml", "")
	flag.Parse()

	err := utils.Execute(*cfg)
	if err != nil {
		panic(err)
	}
}
