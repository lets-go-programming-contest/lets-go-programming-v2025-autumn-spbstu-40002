package main

import (
	"fmt"

	"github.com/XShaygaND/task-8/internal/configs"
)

func main() {
	err := configs.DisplayEnvironmentAndLogLevel()
	if err != nil {
		fmt.Println("Error encountered:", err)
	}
}
