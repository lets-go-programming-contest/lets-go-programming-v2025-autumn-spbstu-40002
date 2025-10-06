package kth

import (
	"container/heap"
	"errors"
)

type minHeap []int

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *minHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func KthMostPreferred(values []int, position int) (int, error) {
	if position < 1 || position > len(values) {
		return 0, errors.New("position out of range")
	}
	h := &minHeap{}
	heap.Init(h)
	for _, v := range values {
		heap.Push(h, v)
		if h.Len() > position {
			heap.Pop(h)
		}
	}
	if h.Len() == 0 {
		return 0, errors.New("empty result")
	}
	return (*h)[0], nil
}
