package main

import (
	"container/heap"
	"fmt"

	"github.com/maryankov.leonid/task-2-2/myheap"
)

func main() {
	var (
		numDish    int
		preference int
	)

	_, err := fmt.Scan(&numDish)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	dishes := make([]int, numDish)

	for index := range dishes {
		_, err = fmt.Scan(&dishes[index])
		if err != nil {
			fmt.Println("Error:", err)

			return
		}
	}

	_, err = fmt.Scan(&preference)
	if err != nil || preference <= 0 {
		fmt.Println("Error:", err)

		return
	}

	dishHeap := myheap.IntHeap(dishes)
	heap.Init(&dishHeap)

	var kDish int

	for range preference {
		val, ok := heap.Pop(&dishHeap).(int)
		if !ok {
			fmt.Println("Error: invalid type")

			return
		}

		kDish = val
	}

	fmt.Println(kDish)
}
