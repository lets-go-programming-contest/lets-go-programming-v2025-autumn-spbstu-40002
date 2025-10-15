package minheap

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x any) {
	value, ok := x.(int)
	if !ok {
		panic("IntHeap.Push: expected int")
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() any {
	// Прочитал, что распространенной практикой является перекладывание ответственности
	// за вызов pop от пустой кучи на вызывающий код, поэтому проверок нет
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}
