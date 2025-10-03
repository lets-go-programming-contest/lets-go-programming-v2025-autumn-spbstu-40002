package main

import (
	"container/heap"
	"fmt"
	"os"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] } // Максимальная куча
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	var dishCount int
	if _, err := fmt.Scan(&dishCount); err != nil {
		fmt.Fprintln(os.Stderr, "couldn't read the number of dishes:", err)
		return
	}

	dishes := make([]int, dishCount)
	for i := 0; i < dishCount; i++ {
		if _, err := fmt.Scan(&dishes[i]); err != nil {
			fmt.Fprintln(os.Stderr, "couldn't read dish preference:", err)
			return
		}
	}

	var k int
	if _, err := fmt.Scan(&k); err != nil {
		fmt.Fprintln(os.Stderr, "couldn't read k:", err)
		return
	}

	h := &IntHeap{}
	heap.Init(h)

	for _, val := range dishes {
		heap.Push(h, val)
	}

	var selected int
	for i := 0; i < k; i++ {
		selected = heap.Pop(h).(int)
	}

	fmt.Println(selected)
}
