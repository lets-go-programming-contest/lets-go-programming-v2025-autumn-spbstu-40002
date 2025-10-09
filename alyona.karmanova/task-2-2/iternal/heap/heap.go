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

func (h *MaxHeap) Push(elem interface{}) {
	val, ok := elem.(int)

	if !ok {
		panic("MaxHeap can only store int values")
	}

	*h = append(*h, val)
}

func (h *MaxHeap) Pop() interface{} {
	if h.Len() == 0 {
		return nil
	}

	oldHeap := *h

	n := len(oldHeap)

	maxElem := oldHeap[n-1]

	*h = oldHeap[0 : n-1]

	return maxElem
}
