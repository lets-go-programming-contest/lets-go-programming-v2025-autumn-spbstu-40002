package interfaces

import (
	"fmt"

	"github.com/slendycs/task-2-2/myerrors"
)

// Define the interface of the heap so that the minimum element is at the top.
type MinHeap []int

func (h *MinHeap) Len() int           { return len(*h) }
func (h *MinHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *MinHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MinHeap) Push(elem any) {
	intElem, success := elem.(int)
	if !success {
		panic(myerrors.ErrIncorrectHeapPushedType)
	}

	*h = append(*h, intElem)
}

func (h *MinHeap) Pop() any {
	if h.Len() == 0 {
		fmt.Println(myerrors.ErrNothingToDelete)

		return nil
	}

	oldHeap := *h
	lenOldHeap := len(oldHeap)
	lastElem := oldHeap[lenOldHeap-1]
	*h = oldHeap[0 : lenOldHeap-1]

	return lastElem
}
