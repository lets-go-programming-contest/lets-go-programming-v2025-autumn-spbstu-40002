package main

import (
	"container/heap"
	"fmt"

	"github.com/eg0sha-0/task-2-2/heaputils"
)

func main() {
	var total int
	if _, err := fmt.Scan(&total); err != nil || total <= 0 {
		return
	}

	arr := make([]int, total)
	for i := range arr {
		if _, err := fmt.Scan(&arr[i]); err != nil {
			return
		}
	}

	var kth int
	if _, err := fmt.Scan(&kth); err != nil || kth <= 0 || kth > total {
		return
	}

	heapInst := &heaputils.IntHeap{}
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
