package kth

import (
	"container/heap"

	"github.com/rekottt/task-2-2/ktherr"
)

type minHeap []int

func (h *minHeap) Len() int           { return len(*h) }
func (h *minHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *minHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *minHeap) Push(x any) {
	v, ok := x.(int)
	if !ok {
		panic("kth: Push expected int value")
	}

	*h = append(*h, v)
}

func (h *minHeap) Pop() any {
	old := *h
	oldLen := len(old)

	if oldLen == 0 {
		panic("kth: Pop called on empty heap")
	}

	x := old[oldLen-1]
	*h = old[:oldLen-1]

	return x
}

func KthMostPreferred(values []int, position int) (int, error) {
	if position < 1 || position > len(values) {
		return 0, ktherr.ErrPositionOutOfRange
	}

	minH := &minHeap{}
	heap.Init(minH)

	for _, v := range values {
		heap.Push(minH, v)

		if minH.Len() > position {
			heap.Pop(minH)
		}
	}

	if minH.Len() == 0 {
		return 0, ktherr.ErrEmptyResult
	}

	return (*minH)[0], nil
}
