package heap

type MaxHeap []int

func (h *MaxHeap) Len() int {
	return len(*h)
}

func (h *MaxHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *MaxHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *MaxHeap) Push(x interface{}) {
	val, ok := x.(int)

	if !ok {
		panic("Value is not integer")
	}

	*h = append(*h, val)
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	length := len(old)

	if length == 0 {
		panic("Heap is empty")
	}

	x := old[length-1]
	*h = old[0 : length-1]

	return x
}
