package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/victor.kim/task-2-2/heaputils"
)

var (
	errInvalidTotal = errors.New("invalid total number of dishes")
	errInvalidInput = errors.New("invalid input while reading dishes")
	errInvalidKth   = errors.New("invalid k value")
)

func main() {
	var total int
	if _, err := fmt.Scan(&total); err != nil || total <= 0 {
		fmt.Println(errInvalidTotal)

		return
	}

	array := make([]int, total)
	for i := range array {
		if _, err := fmt.Scan(&array[i]); err != nil {
			fmt.Println(errInvalidInput)

			return
		}
	}

	var kth int
	if _, err := fmt.Scan(&kth); err != nil || kth <= 0 || kth > total {
		fmt.Println(errInvalidKth)

		return
	}

	minHeap := &heaputils.IntHeap{}
	heap.Init(minHeap)

	for i := range array[:kth] {
		heap.Push(minHeap, array[i])
	}

	for _, value := range array[kth:] {
		if value > (*minHeap)[0] {
			heap.Pop(minHeap)
			heap.Push(minHeap, value)
		}
	}

	fmt.Println((*minHeap)[0])
}
