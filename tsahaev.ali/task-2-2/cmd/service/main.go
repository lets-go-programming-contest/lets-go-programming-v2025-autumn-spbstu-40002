package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	var dishCount, preferenceRank int

	_, err := fmt.Scan(&dishCount)
	if err != nil {
		return
	}

	dishes := make([]int, dishCount)
	for index := range dishes {
		_, err = fmt.Scan(&dishes[index])
		if err != nil {
			return
		}
	}

	_, err = fmt.Scan(&preferenceRank)
	if err != nil {
		return
	}

	result := findKthPreference(dishes, preferenceRank)
	fmt.Println(result)
}

func findKthPreference(dishes []int, preferenceRank int) int {
	heapInstance := &IntHeap{}
	heap.Init(heapInstance)

	for _, dish := range dishes {
		heap.Push(heapInstance, dish)
	}

	for range preferenceRank - 1 {
		heap.Pop(heapInstance)
	}

	return (*heapInstance)[0]
}
