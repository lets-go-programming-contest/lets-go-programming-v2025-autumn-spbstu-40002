package main

import (
	"fmt"

	"github.com/rachguta/task-8/config"
)

func main() {
	cnf := config.GetConfig()
	fmt.Print(cnf.Environment + " " + cnf.Loglevel)
}
