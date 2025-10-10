package main

import (
	"container/heap"
	"fmt"

	workerHeap "github.com/F0LY/task-2-2/internal/heap"
)

func main() {
	var dishCount, preferenceK int
	if _, err := fmt.Scan(&dishCount); err != nil {
		return
	}

	preferences := make([]int, dishCount)
	for i := range preferences {
		if _, err := fmt.Scan(&preferences[i]); err != nil {
			return
		}
	}

	if _, err := fmt.Scan(&preferenceK); err != nil {
		return
	}

	maxHeap := workerHeap.NewMaxHeap()
	for _, rating := range preferences {
		heap.Push(maxHeap, rating)
	}

	var result int

	for range preferenceK {
		popped := heap.Pop(maxHeap)
		if val, ok := popped.(int); ok {
			result = val
		} else {
			panic("heap.Pop returned non-int value")
		}
	}

	fmt.Println(result)
}
