package interfaces

import (
	"fmt"

	"github.com/slendycs/task-2-2/errors"
)

// Define the interface of the heap so that the minimum element is at the top.
type MinHeap []int

func (h *MinHeap) Len() int           { return len(*h) }
func (h *MinHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *MinHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MinHeap) Push(elem any) {
	intElem, ok := elem.(int)
	if !ok {
		fmt.Println(errors.ErrIncorrectHeapPushedType)
		
		return
	}
	*h = append(*h, intElem)
}

func (h *MinHeap) Pop() any {
	oldHeap := *h
	lenOldHeap := len(oldHeap)
	lastElem := oldHeap[lenOldHeap-1]
	*h = oldHeap[0 : lenOldHeap-1]

	return lastElem
}
