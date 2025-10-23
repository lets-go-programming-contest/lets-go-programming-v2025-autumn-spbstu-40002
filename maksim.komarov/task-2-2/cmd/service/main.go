package main

import (
	"container/heap"
	"fmt"

	"github.com/megurumacabre/task-2-2/internal/maxheap"
)

func main() {
	var dishCount int
	if _, err := fmt.Scan(&dishCount); err != nil {
		fmt.Println("invalid dishes count")

		return
	}

	ratings := make([]int, dishCount)

	for i := range dishCount {
		if _, err := fmt.Scan(&ratings[i]); err != nil {
			fmt.Println("invalid rating")

			return
		}
	}

	var kthIndex int
	if _, err := fmt.Scan(&kthIndex); err != nil {
		fmt.Println("invalid k")

		return
	}

	dishHeap := &maxheap.MaxHeap{}
	heap.Init(dishHeap)

	for i := range dishCount {
		heap.Push(dishHeap, ratings[i])
	}

	selected := 0

	for range kthIndex {
		val, ok := heap.Pop(dishHeap).(int)
		if !ok {
			fmt.Println("heap pop type error")

			return
		}

		selected = val
	}

	fmt.Println(selected)
}
