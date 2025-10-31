package interheap

import "container/heap"

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
	value, ok := x.(int)

	if !ok {
		panic("IntHeap.Push: value is not int")
	}

	h.data = append(h.data, value)
}

func (h *IntHeap) Pop() interface{} {
	if len(h.data) == 0 {
		panic("IntHeap.Pop: empty IntHeap")
	}

	n := len(h.data)
	x := h.data[n-1]
	h.data = h.data[:n-1]

	return x
}

func (h *IntHeap) Peek() int {
	if len(h.data) == 0 {
		panic("empty")
	}

	return h.data[0]
}

func FindKthLargest(nums []int, kth int) int {
	meanHeap := &IntHeap{
		data: []int{},
	}
	heap.Init(meanHeap)

	for _, num := range nums {
		heap.Push(meanHeap, num)

		if meanHeap.Len() > kth {
			heap.Pop(meanHeap)
		}
	}

	return meanHeap.Peek()
}
