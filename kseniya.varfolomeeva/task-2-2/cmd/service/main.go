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
	var dishCount int
	fmt.Scan(&dishCount)

	dishes := make([]int, dishCount)
	for i := range dishes {
		fmt.Scan(&dishes[i])
	}

	var preferenceOrder int
	fmt.Scan(&preferenceOrder)

	heapInstance := &MaxHeap{}
	heap.Init(heapInstance)

	for _, num := range dishes {
		heap.Push(heapInstance, num)
	}

	var result int
	for i := 0; i < preferenceOrder; i++ {
		result = heap.Pop(heapInstance).(int)
	}

	fmt.Println(result)
}
