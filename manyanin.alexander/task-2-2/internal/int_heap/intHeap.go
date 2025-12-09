package intheap

import "errors"

var (
	ErrIncorrectType = errors.New("incorrect type pushed to IntHeap")
	ErrEmptyHeap     = errors.New("pop from empty heap")
)

type IntHeap []int

func (heap *IntHeap) Len() int {
	return len(*heap)
}

func (heap *IntHeap) Less(firstIndex, secondIndex int) bool {
	return (*heap)[firstIndex] > (*heap)[secondIndex] // max-heap получается; если < - min-heap
}

func (heap *IntHeap) Swap(firstIndex, secondIndex int) {
	(*heap)[firstIndex], (*heap)[secondIndex] = (*heap)[secondIndex], (*heap)[firstIndex]
}

func (heap *IntHeap) Push(value interface{}) {
	convertedValue, correctType := value.(int)
	if !correctType {
		panic(ErrIncorrectType)
	}

	*heap = append(*heap, convertedValue)
}

func (heap *IntHeap) Pop() interface{} {
	oldHeap := *heap

	heapLength := len(oldHeap)
	if heapLength == 0 {
		return nil
	}

	value := oldHeap[heapLength-1]

	*heap = oldHeap[0 : heapLength-1]

	return value
}
