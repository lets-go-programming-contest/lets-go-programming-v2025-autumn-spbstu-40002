package heap

import "container/heap"

type MinHeap []int

func (h *MinHeap) Len() int {
	return len(*h)
}

func (h *MinHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *MinHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *MinHeap) Push(x any) {
	value, ok := x.(int)

	if !ok {
		panic("heap.Push: value is not int")
	}

	*h = append(*h, value)
}

func (h *MinHeap) Pop() any {
	old := *h
	length := len(old)

	if length == 0 {
		panic("heap.Pop: empty heap")
	}

	x := old[length-1]
	*h = old[:length-1]

	return x
}

func FindKthLargest(dishes []int, kth int) int {
	minHeap := &MinHeap{}

	heap.Init(minHeap)

	for _, dish := range dishes {
		heap.Push(minHeap, dish)

		if minHeap.Len() > kth {
			heap.Pop(minHeap)
		}
	}

	return (*minHeap)[0]
}
