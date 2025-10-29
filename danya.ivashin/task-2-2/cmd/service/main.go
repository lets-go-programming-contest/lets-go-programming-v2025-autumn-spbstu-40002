package main

import (
	"container/heap"
	"fmt"

	"github.com/Danya-byte/task-2-2/internal/intheap"
)

func getKDish(dishes []int, k int) int {
	window := new(intheap.IntHeap)
	*window = intheap.IntHeap(dishes[:k])
	heap.Init(window)

	dishes = dishes[k:]
	for _, cost := range dishes {
		top, _ := window.Top()
		if cost > top {
			window.ReplaceTop(cost)
		}
	}

	top, _ := window.Top()

	return top
}

func main() {
	var dishesCount int
	if _, err := fmt.Scan(&dishesCount); err != nil || dishesCount <= 0 {
		fmt.Printf("failed to dishes count: %s\n", err)

		return
	}

	dishes := make([]int, dishesCount)

	for index := range dishesCount {
		var cost int
		if _, err := fmt.Scan(&cost); err != nil {
			fmt.Printf("failed to scan cost: %s\n", err)

			return
		}

		dishes[index] = cost
	}

	var need int
	if _, err := fmt.Scan(&need); err != nil || need <= 0 {
		fmt.Printf("failed to scan needed element: %s\n", err)

		return
	}

	if need > dishesCount {
		fmt.Println("dishes count should be less then needed element")

		return
	}

	fmt.Println(getKDish(dishes, need))
}
