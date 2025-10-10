package maxheap

type MaxHeap []int

func (h *MaxHeap) Len() int           { return len(*h) }
func (h *MaxHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *MaxHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MaxHeap) Push(x any) {
	v, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, v)
}

func (h *MaxHeap) Pop() any {
	old := *h
	if len(old) == 0 {
		return 0
	}

	v := old[len(old)-1]
	*h = old[:len(old)-1]

	return v
}
