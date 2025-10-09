package main

import (
	"container/heap"
	"fmt"
	"os"
)

type IntHeap []int

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x interface{}) {
	v, ok := x.(int)
	if !ok {
		return
	}
	*h = append(*h, v)
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]

	return x
}

func main() {
	var count int
	if _, err := fmt.Fscan(os.Stdin, &count); err != nil || count <= 0 {
		return
	}

	values := make([]int, count)
	for i := range make([]struct{}, count) {
		if _, err := fmt.Fscan(os.Stdin, &values[i]); err != nil {
			return
		}
	}

	var kth int
	if _, err := fmt.Fscan(os.Stdin, &kth); err != nil || kth <= 0 || kth > count {
		return
	}

	minHeap := &IntHeap{}
	heap.Init(minHeap)

	for i := range make([]struct{}, kth) {
		heap.Push(minHeap, values[i])
	}

	for j := range make([]struct{}, count-kth) {
		idx := kth + j
		val := values[idx]
		if val > (*minHeap)[0] {
			heap.Pop(minHeap)
			heap.Push(minHeap, val)
		}
	}

	top := (*minHeap)[0]
	fmt.Println(top)
}
