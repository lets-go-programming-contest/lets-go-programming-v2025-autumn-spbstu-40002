package main

import (
	"fmt"

	"github.com/Expeline/task-2-2/internal/heap"
)

const (
	MinValue = -10000
	MaxValue = 10000
	MinN     = 1
	MaxN     = 10000
)

func main() {
	var dishesCount int

	if _, err := fmt.Scan(&dishesCount); err != nil {
		return
	}

	if dishesCount < MinN || dishesCount > MaxN {
		return
	}

	dishes := []int{}

	for range dishesCount {
		var value int

		if _, err := fmt.Scan(&value); err != nil {
			return
		}

		if value < MinValue || value > MaxValue {
			return
		}

		dishes = append(dishes, value)
	}

	var kthDish int

	if _, err := fmt.Scan(&kthDish); err != nil {
		return
	}

	if kthDish < 1 || kthDish > dishesCount {
		return
	}

	result := heap.FindKthLargest(dishes, kthDish)

	fmt.Println(result)
}
