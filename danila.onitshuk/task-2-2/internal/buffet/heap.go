package buffet

type Heap []int

func (heap *Heap) Len() int           { return len(*heap) }
func (heap *Heap) Less(i, j int) bool { return (*heap)[i] > (*heap)[j] }
func (heap *Heap) Swap(i, j int)      { (*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i] }

func (heap *Heap) Push(element any) {
	intElement, err := element.(int)
	if !err {
		return
	}

	*heap = append(*heap, intElement)
}

func (heap *Heap) Pop() any {
	if heap.Len() == 0 {
		return ErrUnderFlow
	}

	oldHeap := *heap
	newHeapSize := len(oldHeap) - 1
	lastElement := oldHeap[newHeapSize]
	*heap = oldHeap[0:newHeapSize]

	return lastElement
}
