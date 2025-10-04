package heap

type Heap []int

func (heap Heap) Len() int           { return len(heap) }
func (heap Heap) Less(i, j int) bool { return heap[i] > heap[j] }
func (heap Heap) Swap(i, j int)      { heap[i], heap[j] = heap[j], heap[i] }

func (heap *Heap) Push(newHeap interface{}) {
	*heap = append(*heap, newHeap.(int))
}

func (heap *Heap) Pop() interface{} {
	oldHeap := *heap
	newHeapSize := len(oldHeap) - 1
	lastElement := oldHeap[newHeapSize]
	*heap = oldHeap[0:newHeapSize]
	return lastElement
}
