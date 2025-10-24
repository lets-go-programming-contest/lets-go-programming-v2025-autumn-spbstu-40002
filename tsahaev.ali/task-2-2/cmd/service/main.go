package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	var N, k int
	fmt.Scan(&N)

	dishes := make([]int, N)
	for i := 0; i < N; i++ {
		fmt.Scan(&dishes[i])
	}

	fmt.Scan(&k)

	result := findKthPreference(dishes, k)
	fmt.Println(result)
}

func findKthPreference(dishes []int, k int) int {
	h := &IntHeap{}
	heap.Init(h)

	for _, dish := range dishes {
		heap.Push(h, dish)
	}

	for i := 0; i < k-1; i++ {
		heap.Pop(h)
	}

	return (*h)[0]
}
