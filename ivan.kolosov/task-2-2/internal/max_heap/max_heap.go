package maxheap

type MaxHeap []int

func (heap *MaxHeap) Len() int { return len((*heap)) }

func (heap *MaxHeap) Less(i, j int) bool { return (*heap)[i] > (*heap)[j] }

func (heap *MaxHeap) Swap(i, j int) { (*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i] }

func (heap *MaxHeap) Push(x interface{}) {
	value, isInt := x.(int)
	if !isInt {
		panic("expected int")
	}

	*heap = append(*heap, value)
}

func (heap *MaxHeap) Pop() interface{} {
	if len(*heap) == 0 {
		return nil
	}

	x := (*heap)[len(*heap)-1]
	*heap = (*heap)[0 : len(*heap)-1]

	return x
}
