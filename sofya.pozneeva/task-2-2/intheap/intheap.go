package intheap

import (
	"fmt"
)

type IntHeap []int16

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x any) {
	if val, ok := x.(int16); ok {
		*h = append(*h, val)
	} else {
		panic(fmt.Sprintf("expected int16, got %T", x))
	}
}

func (h *IntHeap) Pop() any {
	if len(*h) != 0 {
		old := *h
		n := len(old)
		x := old[n-1]
		*h = old[0 : n-1]

		return x
	} else {
		panic("IntHeap.Pop: heap length must be greater than zero")
	}
}
