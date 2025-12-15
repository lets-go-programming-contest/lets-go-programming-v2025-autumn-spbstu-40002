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
	v := x.(int)
	*h = append(*h, v)
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	var n, k int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}

	h := &IntHeap{}
	heap.Init(h)

	for i := 0; i < n; i++ {
		var x int
		if _, err := fmt.Scan(&x); err != nil {
			return
		}
		heap.Push(h, x)
	}

	if _, err := fmt.Scan(&k); err != nil {
		return
	}

	for i := 1; i < k; i++ {
		heap.Pop(h)
	}

	fmt.Println(heap.Pop(h).(int))
}
