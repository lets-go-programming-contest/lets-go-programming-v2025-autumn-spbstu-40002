package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type MaxHeap []int

func (h *MaxHeap) Len() int           { return len(*h) }
func (h *MaxHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *MaxHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MaxHeap) Push(v interface{}) {
	val, ok := v.(int)
	if !ok {
		return
	}
	*h = append(*h, val)
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]

	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var count int
	if _, err := fmt.Fscan(reader, &count); err != nil {
		os.Exit(1)
	}

	if count <= 0 {
		os.Exit(1)
	}

	maxHeap := &MaxHeap{}
	heap.Init(maxHeap)

	for i := 0; i < count; i++ {
		var x int
		if _, err := fmt.Fscan(reader, &x); err != nil {
			os.Exit(1)
		}

		heap.Push(maxHeap, x)
	}

	var order int
	if _, err := fmt.Fscan(reader, &order); err != nil {
		os.Exit(1)
	}

	if order < 1 || order > count {
		os.Exit(1)
	}

	var result int
	for i := 0; i < order; i++ {
		if v, ok := heap.Pop(maxHeap).(int); ok {
			result = v
		} else {
			os.Exit(1)
		}
	}

	fmt.Println(result)
}

