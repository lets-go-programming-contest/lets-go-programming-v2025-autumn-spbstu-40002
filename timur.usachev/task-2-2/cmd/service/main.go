package main

import (
	"container/heap"
	"fmt"
	"os"
)

type IntHeap struct {
	data []int
}

func (h *IntHeap) Len() int {
	return len(h.data)
}

func (h *IntHeap) Less(i, j int) bool {
	return h.data[i] < h.data[j]
}

func (h *IntHeap) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *IntHeap) Push(x interface{}) {
	v, ok := x.(int)
	if !ok {
		return
	}

	h.data = append(h.data, v)
}

func (h *IntHeap) Pop() interface{} {
	if len(h.data) == 0 {
		return nil
	}

	n := len(h.data)
	x := h.data[n-1]
	h.data = h.data[:n-1]

	return x
}

func main() {
	var count int
	if _, err := fmt.Fscan(os.Stdin, &count); err != nil || count <= 0 {
		fmt.Fprintln(os.Stderr, "invalid input: count must be a positive integer")
		os.Exit(1)
	}

	values := make([]int, count)
	for i := range make([]struct{}, count) {
		if _, err := fmt.Fscan(os.Stdin, &values[i]); err != nil {
			fmt.Fprintln(os.Stderr, "invalid input: failed to read values")
			os.Exit(1)
		}
	}

	var kth int
	if _, err := fmt.Fscan(os.Stdin, &kth); err != nil || kth <= 0 || kth > count {
		fmt.Fprintln(os.Stderr, "invalid input: k must be between 1 and count")
		os.Exit(1)
	}

	minHeap := &IntHeap{data: []int{}}
	heap.Init(minHeap)

	for i := range make([]struct{}, kth) {
		heap.Push(minHeap, values[i])
	}

	for j := range make([]struct{}, count-kth) {
		val := values[kth+j]

		if minHeap.Len() == 0 {
			continue
		}

		if val > minHeap.data[0] {
			heap.Pop(minHeap)
			heap.Push(minHeap, val)
		}
	}

	if minHeap.Len() == 0 {
		fmt.Fprintln(os.Stderr, "no result")
		os.Exit(1)
	}

	top := minHeap.data[0]
	fmt.Println(top)
}
