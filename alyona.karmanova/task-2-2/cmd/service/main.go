package main

import (
	"container/heap"
	"fmt"

	myheap "github.com/HuaChenju/task-2-2/iternal/heap"
)

const (
	MinCount      = 1
	MaxCount      = 10000
	MinRating     = -10000
	MaxRating     = 10000
	MinPreference = 1
)

func main() {
	var count, preference, rating int

	heapMax := &myheap.MaxHeap{}

	_, err := fmt.Scan(&count)

	if err != nil || count < MinCount || count > MaxCount {
		fmt.Println("incorrect amount of dishes")

		return
	}

	for range count {
		_, err = fmt.Scan(&rating)

		if err != nil || rating < MinRating || rating > MaxRating {
			fmt.Println("incorrect rating")
			return
		}

		heap.Push(heapMax, rating)
	}

	_, err = fmt.Scan(&preference)

	if err != nil || preference > count || preference < MinPreference {
		return
	}

	for range preference - 1 {
		val := heap.Pop(heapMax)
		if val == nil {
			fmt.Println("Heap is empty unexpectedly")

			return
		}
	}

	val := heap.Pop(heapMax)
	if val == nil {
		fmt.Println("Heap is empty unexpectedly")

		return
	}

	fmt.Println(val)
}
