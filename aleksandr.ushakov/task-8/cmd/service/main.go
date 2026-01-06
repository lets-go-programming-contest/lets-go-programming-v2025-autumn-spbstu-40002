package main

import (
	"fmt"

	"github.com/rachguta/task-8/config"
)

func main() {
	fmt.Println(config.Cnf.Environment + " " + config.Cnf.Log_level)
}
