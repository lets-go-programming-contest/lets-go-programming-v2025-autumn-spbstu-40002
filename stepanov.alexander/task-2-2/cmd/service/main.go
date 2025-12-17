package main

import (
	"container/heap"
	"fmt"
)

type IntHeap struct {
	data []int
}

func (h *IntHeap) Len() int {
	return len(h.data)
}

func (h *IntHeap) Less(i, j int) bool {
	return h.data[i] > h.data[j]
}

func (h *IntHeap) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *IntHeap) Push(x any) {
	value, ok := x.(int)
	if !ok {
		return
	}

	h.data = append(h.data, value)
}

func (h *IntHeap) Pop() any {
	n := len(h.data)
	value := h.data[n-1]
	h.data = h.data[:n-1]

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

	priorityQueue := &IntHeap{data: ratings}
	heap.Init(priorityQueue)

	if _, err := fmt.Scan(&order); err != nil {
		return
	}

	for i := 1; i < order; i++ {
		heap.Pop(priorityQueue)
	}

	result, ok := heap.Pop(priorityQueue).(int)
	if !ok {
		return
	}

	fmt.Println(result)
}
