package main

import (
	"container/heap"
	"fmt"
)

type MaxHeap []int

func (h *MaxHeap) Len() int {
	return len(*h)
}

func (h *MaxHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *MaxHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *MaxHeap) Push(elem interface{}) {
	val, err := elem.(int)

	if !err {
		return
	}

	*h = append(*h, val)
}

func (h *MaxHeap) Pop() interface{} {
	oldHeap := *h

	n := len(oldHeap)

	maxElem := oldHeap[n-1]

	*h = oldHeap[0 : n-1]

	return maxElem
}

func main() {
	var count, preference, rating int

	heapMax := &MaxHeap{}

	_, err := fmt.Scan(&count)

	if err != nil || count < 1 || count > 10000 {
		return
	}

	for range count {
		_, err = fmt.Scan(&rating)

		if err != nil || rating < -10000 || rating > 10000 {
			return
		}

		heap.Push(heapMax, rating)
	}

	_, err = fmt.Scan(&preference)

	if err != nil || preference > count {
		return
	}

	for range preference - 1 {
		heap.Pop(heapMax)
	}

	fmt.Println(heap.Pop(heapMax))
}
