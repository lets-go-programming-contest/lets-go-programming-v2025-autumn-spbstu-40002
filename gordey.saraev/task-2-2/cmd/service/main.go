package main

import (
	"container/heap"
	"fmt"
)

type MaxHeap []int

func (h *MaxHeap) Len() int           { return len(*h) }
func (h *MaxHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *MaxHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MaxHeap) Push(value interface{}) {
	if v, ok := value.(int); ok {
		*h = append(*h, v)
	}
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]

	return x
}

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

	maxHeap := &MaxHeap{}
	heap.Init(maxHeap)

	for _, rating := range preferences {
		heap.Push(maxHeap, rating)
	}

	var result int

	for range preferenceK {
		popped := heap.Pop(maxHeap)
		if val, ok := popped.(int); ok {
			result = val
		}
	}

	fmt.Println(result)
}
