package main

import (
	"fmt"
	"github.com/bolatbyek/task-2-1/internal/department"
)

func main() {
	var N, K int
	fmt.Scan(&N, &K)

	for i := 0; i < N; i++ {
		// Process each department
		manager := department.NewManager()

		for j := 0; j < K; j++ {
			var constraint string
			fmt.Scan(&constraint)

			// Add constraint and get optimal temperature
			optimalTemp := manager.AddConstraint(constraint)
			fmt.Println(optimalTemp)
		}
	}
}