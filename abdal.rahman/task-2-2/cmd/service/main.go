package main

import (
	"container/heap"
	"fmt"

	"github.com/xkoex/task-2-2/myheap"
)

func main() {
	var numDish int

	_, err := fmt.Scan(&numDish)
	if err != nil || numDish <= 0 {

		return
	}

	dishes := make([]int, numDish)

	for i := range dishes {
		_, err = fmt.Scan(&dishes[i])
		if err != nil {

			return
		}
	}

	var rankChoice int

	_, err = fmt.Scan(&rankChoice)
	if err != nil || rankChoice <= 0 {

		return
	}

	dishHeap := myheap.NumHeap(dishes)
	heap.Init(&dishHeap)

	var selectedDish int

	for i := 0; i < rankChoice; i++ {
		val, ok := heap.Pop(&dishHeap).(int)
		if !ok {

			return
		}

		selectedDish = val
	}

	fmt.Println(selectedDish)
}
