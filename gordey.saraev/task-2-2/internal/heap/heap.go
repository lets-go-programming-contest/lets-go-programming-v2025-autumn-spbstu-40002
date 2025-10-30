package heap

import "container/heap"

type MaxHeap []int

func (h *MaxHeap) Len() int           { return len(*h) }
func (h *MaxHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *MaxHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MaxHeap) Push(value interface{}) {
	v, ok := value.(int)
	if !ok {
		panic("MaxHeap.Push: value is not int")
	}

	*h = append(*h, v)
}

func (h *MaxHeap) Pop() interface{} {
	if len(*h) == 0 {
		panic("MaxHeap.Pop: heap is empty")
	}

	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]

	return x
}

func NewMaxHeap() *MaxHeap {
	h := &MaxHeap{}
	heap.Init(h)

	return h
}
