package intheap

import "container/heap"

type IntHeap []int

func (h *IntHeap) Push(x interface{}) {
	if val, ok := x.(int); ok {
		*h = append(*h, val)
	} else {
		return
	}
}

func (h *IntHeap) Pop() interface{} {
	intHeap := *h
	x := intHeap[len(intHeap)-1]
	*h = intHeap[:len(intHeap)-1]

	return x
}

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) PushValue(x int) {
	heap.Push(h, x)
}

func (h *IntHeap) PopValue() int {
	if val, ok := h.Pop().(int); ok {
		return val
	} else {
		return -1
	}
}
