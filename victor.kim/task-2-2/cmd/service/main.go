package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/victor.kim/task-2-2/heaputils"
)

var (
	errInvalidTotal = errors.New("invalid number of elements")
	errInvalidValue = errors.New("invalid input value")
	errInvalidK     = errors.New("invalid k value")
)

func main() {
	var total int
	if _, err := fmt.Scan(&total); err != nil || total <= 0 {
		fmt.Println(errInvalidTotal)
		return
	}

	arr := make([]int, total)
	for i := range arr {
		if _, err := fmt.Scan(&arr[i]); err != nil {
			fmt.Println(errInvalidValue)
			return
		}
	}

	var kth int
	if _, err := fmt.Scan(&kth); err != nil || kth <= 0 || kth > total {
		fmt.Println(errInvalidK)
		return
	}

	h := &heaputils.IntHeap{}
	heap.Init(h)

	for i := range arr[:kth] {
		heap.Push(h, arr[i])
	}

	for _, v := range arr[kth:] {
		if v > (*h)[0] {
			heap.Pop(h)
			heap.Push(h, v)
		}
	}

	fmt.Println((*h)[0])
}
