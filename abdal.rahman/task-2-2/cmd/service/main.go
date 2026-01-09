package main

import (
	"container/heap"
	"fmt"

	"github.com/xkoex/task-2-2/myheap"
)

func main() {
	var totalDishes int
	var rankChoice int

	_, err := fmt.Scan(&totalDishes)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	dishValues := make([]int, totalDishes)
	for i := range dishValues {
		_, err = fmt.Scan(&dishValues[i])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	_, err = fmt.Scan(&rankChoice)
	if err != nil || rankChoice <= 0 {
		fmt.Println("Error:", err)
		return
	}

	priorityHeap := myheap.NumHeap(dishValues)
	heap.Init(&priorityHeap)

	var selectedDish int
	for i := 0; i < rankChoice; i++ {
		val, ok := heap.Pop(&priorityHeap).(int)
		if !ok {
			fmt.Println("Error: invalid type")
			return
		}
		selectedDish = val
	}

	fmt.Println(selectedDish)
}
