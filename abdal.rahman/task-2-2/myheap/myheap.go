package myheap

type NumHeap []int

func (h *NumHeap) Len() int {
	return len(*h)
}

func (h *NumHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *NumHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *NumHeap) Push(x any) {
	num, ok := x.(int)
	if !ok {
		panic("NumHeap: Push value is not an int")
	}
	*h = append(*h, num)
}

func (h *NumHeap) Pop() any {
	oldHeap := *h
	n := len(oldHeap)
	if n == 0 {
		panic("NumHeap: Pop from empty heap")
	}
	val := oldHeap[n-1]
	*h = oldHeap[0 : n-1]
	return val
}
