package main

import (
	"container/heap"
	"fmt"

	"github.com/xkoex/task-2-2/myheap"
)

func main() {
	var numDish int
	var choiceRank int

	_, err := fmt.Scan(&numDish)
	if err != nil || numDish <= 0 {
		fmt.Println("Error:", err)
		return
	}

	dishes := make([]int, numDish)
	for i := range dishes {
		_, err := fmt.Scan(&dishes[i])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	_, err = fmt.Scan(&choiceRank)
	if err != nil || choiceRank <= 0 {
		fmt.Println("Error:", err)
		return
	}

	platesHeap := myheap.NumHeap(dishes)
	heap.Init(&platesHeap)

	var selectedPlate int
	for range make([]struct{}, choiceRank) {
		val, ok := heap.Pop(&platesHeap).(int)
		if !ok {
			fmt.Println("Error: invalid type")
			return
		}
		selectedPlate = val
	}

	fmt.Println(selectedPlate)
}
