package heaputils

type IntHeap []int

func (ih *IntHeap) Len() int {
	return len(*ih)
}

func (ih *IntHeap) Less(i, j int) bool {
	return (*ih)[i] < (*ih)[j] // min-heap
}

func (ih *IntHeap) Swap(i, j int) {
	(*ih)[i], (*ih)[j] = (*ih)[j], (*ih)[i]
}

func (ih *IntHeap) Push(x interface{}) {
	v, ok := x.(int)
	if !ok {
		panic("IntHeap: Push expects int")
	}

	*ih = append(*ih, v)
}

func (ih *IntHeap) Pop() interface{} {
	old := *ih
	if len(old) == 0 {
		return nil
	}

	lastIndex := len(old) - 1
	last := old[lastIndex]
	*ih = old[:lastIndex]

	return last
}
