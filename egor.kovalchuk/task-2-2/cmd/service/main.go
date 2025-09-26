package main

import (
	"container/heap"
	"fmt"
)

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
	n := len(old)
	x := old[n-1]
	*ih = old[:n-1]

	return x
}

func main() {
	var total int
	if _, err := fmt.Scan(&total); err != nil {
		return
	}

	arr := make([]int, total)
	for i := range arr {
		if _, err := fmt.Scan(&arr[i]); err != nil {
			return
		}
	}

	var kth int
	if _, err := fmt.Scan(&kth); err != nil {
		return
	}

	if kth <= 0 || kth > total {
		return
	}

	heapInst := &IntHeap{}
	heap.Init(heapInst)

	for i := range arr[:kth] {
		heap.Push(heapInst, arr[i])
	}

	for _, v := range arr[kth:] {
		if v > (*heapInst)[0] {
			heap.Pop(heapInst)
			heap.Push(heapInst, v)
		}
	}

	fmt.Println((*heapInst)[0])
}
