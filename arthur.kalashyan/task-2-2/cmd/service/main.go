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

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *MinHeap) Pop() any {
	old := *h

	n := len(old)

	x := old[n-1]

	*h = old[:n-1]

	return x
}

func findKthLargest(dishes []int, k int) int {
	h := &MinHeap{}

	heap.Init(h)

	for _, dish := range dishes {
		heap.Push(h, dish)

		if h.Len() > k {
			heap.Pop(h)
		}
	}

	return (*h)[0]
}

func main() {
	var count int

	if _, err := fmt.Scan(&count); err != nil {
		return
	}

	if count < MinN || count > MaxN {
		return
	}

	dishes := []int{}

	for range count {
		var value int

		if _, err := fmt.Scan(&value); err != nil {
			return
		}

		if value < MinValue || value > MaxValue {
			return
		}

		dishes = append(dishes, value)
	}

	var k int

	if _, err := fmt.Scan(&k); err != nil {
		return
	}

	if k < 1 || k > count {
		return
	}

	result := findKthLargest(dishes, k)

	fmt.Println(result)
}
