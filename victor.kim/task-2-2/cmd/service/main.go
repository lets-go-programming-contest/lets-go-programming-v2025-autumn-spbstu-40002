package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		panic("IntHeap: Push expects int")
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() interface{} {
	oldHeap := *h
	heapSize := len(oldHeap)

	if heapSize == 0 {
		return nil
	}

	lastElement := oldHeap[heapSize-1]
	*h = oldHeap[:heapSize-1]

	return lastElement
}

func main() {
	var total int
	if _, err := fmt.Scan(&total); err != nil || total <= 0 {
		return
	}

	array := make([]int, total)
	for index := range array {
		if _, err := fmt.Scan(&array[index]); err != nil {
			return
		}
	}

	var kth int
	if _, err := fmt.Scan(&kth); err != nil || kth <= 0 || kth > total {
		return
	}

	minHeap := &IntHeap{}
	heap.Init(minHeap)

	for index := range array[:kth] {
		heap.Push(minHeap, array[index])
	}

	for _, value := range array[kth:] {
		if value > (*minHeap)[0] {
			heap.Pop(minHeap)
			heap.Push(minHeap, value)
		}
	}

	fmt.Println((*minHeap)[0])
}
