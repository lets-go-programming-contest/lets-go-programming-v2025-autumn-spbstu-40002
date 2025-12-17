package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] > h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x any) {
	value := x.(int)
	*h = append(*h, value)
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	value := old[n-1]
	*h = old[:n-1]

	return value
}

func main() {
	var dishesCount, order int

	if _, err := fmt.Scan(&dishesCount); err != nil {
		return
	}

	ratings := make([]int, dishesCount)
	for i := range ratings {
		if _, err := fmt.Scan(&ratings[i]); err != nil {
			return
		}
	}

	h := IntHeap(ratings)
	heap.Init(&h)

	if _, err := fmt.Scan(&order); err != nil {
		return
	}

	for i := 1; i < order; i++ {
		heap.Pop(&h)
	}

	fmt.Println(heap.Pop(&h).(int))
}
