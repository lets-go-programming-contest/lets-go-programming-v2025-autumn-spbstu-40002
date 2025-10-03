package main

import (
	"container/heap"
	"fmt"
)

const (
	MinValue = -10000
	MaxValue = 10000
	MinN     = 1
	MaxN     = 10000
)

type MinHeap []int

func (h *MinHeap) Len() int {
	return len(*h)
}

func (h *MinHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *MinHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *MinHeap) Push(x any) {
	value, ok := x.(int)

	if !ok {
		return
	}

	*h = append(*h, value)
}

func (h *MinHeap) Pop() any {
	old := *h

	n := len(old)

	x := old[n-1]

	*h = old[:n-1]

	return x
}

func findKthLargest(dishes []int, kth int) int {
	minHeap := &MinHeap{}

	heap.Init(minHeap)

	for _, dish := range dishes {
		heap.Push(minHeap, dish)

		if minHeap.Len() > kth {
			heap.Pop(minHeap)
		}
	}

	return (*minHeap)[0]
}

func main() {
	var dishesCount int

	if _, err := fmt.Scan(&dishesCount); err != nil {
		return
	}

	if dishesCount < MinN || dishesCount > MaxN {
		return
	}

	dishes := []int{}

	for i := 0; i < dishesCount; i++ {
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

	result := findKthLargest(dishes, kthDish)

	fmt.Println(result)
}
