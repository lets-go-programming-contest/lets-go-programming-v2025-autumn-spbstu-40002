package main

import (
	"container/heap"
	"fmt"
)

type MaxHeap []int

func (h *MaxHeap) Len() int           { return len(*h) }
func (h *MaxHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *MaxHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

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
	var dishCount, preferenceK int
	fmt.Scan(&dishCount)

	preferences := make([]int, dishCount)
	for i := 0; i < dishCount; i++ {
		fmt.Scan(&preferences[i])
	}
	fmt.Scan(&preferenceK)

	h := &MaxHeap{}
	heap.Init(h)
	for _, x := range preferences {
		heap.Push(h, x)
	}

	var result int
	for i := 0; i < preferenceK; i++ {
		result = heap.Pop(h).(int)
	}

	fmt.Println(result)
}
