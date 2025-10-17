package main

import (
	"container/heap"
	"fmt"
)

type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	var count int
	_, err := fmt.Scan(&count)
	if err != nil {
		return
	}

	dishes := make([]int, count)
	for i := 0; i < count; i++ {
		_, err = fmt.Scan(&dishes[i])
		if err != nil {
			return
		}
	}

	var topCount int
	_, err = fmt.Scan(&topCount)
	if err != nil {
		return
	}

	maxHeap := &MaxHeap{}
	heap.Init(maxHeap)

	for _, dish := range dishes {
		heap.Push(maxHeap, dish)
	}

	var result int
	for i := 0; i < topCount; i++ {
		value := heap.Pop(maxHeap)
		result = value.(int)
	}

	fmt.Println(result)
}