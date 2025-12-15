package intheap

import "container/heap"

type IntHeap []int

func (h *IntHeap) Push(x interface{}) {
	if val, ok := x.(int); ok {
		*h = append(*h, val)
	} else {
		panic("IntHeap.Push: expected int type")
	}
}

func (h *IntHeap) Pop() interface{} {
	if h.Len() == 0 {
		panic("IntHeap.Pop: heap length must be greater than zero")
	}

	intHeap := *h
	x := intHeap[len(intHeap)-1]
	*h = intHeap[:len(intHeap)-1]
	heap.Fix(h, 0)

	return x
}

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}
