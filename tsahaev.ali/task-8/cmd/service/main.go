package main

import (
	"fmt"

	"github.com/Tsahaev/task-8/internal/configs"
)

func main() {
	err := configs.PrintTagAndLevel()
	if err != nil {
		fmt.Println("Error output:", err)
	}
}
