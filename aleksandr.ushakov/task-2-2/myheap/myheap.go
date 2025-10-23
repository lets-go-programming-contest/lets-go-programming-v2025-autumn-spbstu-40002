package myheap

// Heap realization.
type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x any) {
	val, ok := x.(int)
	if !ok {
		panic("can not convert interface to int")
	}

	*h = append(*h, val)
}

func (h *IntHeap) Pop() any {
	old := *h
	length := len(old)

	if length == 0 {
		panic("heap is empty")
	}

	x := old[length-1]
	*h = old[0 : length-1]

	return x
}
