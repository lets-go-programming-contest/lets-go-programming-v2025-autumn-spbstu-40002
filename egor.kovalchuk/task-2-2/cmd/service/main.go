package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] } // min-heap
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	var N int
	if _, err := fmt.Scan(&N); err != nil {
		return
	}

	arr := make([]int, N)
	for i := 0; i < N; i++ {
		if _, err := fmt.Scan(&arr[i]); err != nil {
			return
		}
	}

	var k int
	if _, err := fmt.Scan(&k); err != nil {
		return
	}

	h := &IntHeap{}
	heap.Init(h)

	for i := 0; i < k; i++ {
		heap.Push(h, arr[i])
	}

	for i := k; i < N; i++ {
		if arr[i] > (*h)[0] {
			heap.Pop(h)
			heap.Push(h, arr[i])
		}
	}

	fmt.Println((*h)[0])
}
