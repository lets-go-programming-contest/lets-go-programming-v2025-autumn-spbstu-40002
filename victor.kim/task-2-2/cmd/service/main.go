package main

import (
	"container/heap"
	"fmt"
	"os"
)

type DishHeap []int

func (h *DishHeap) Len() int {
	return len(*h)
}

func (h *DishHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j] // максимальная куча
}

func (h *DishHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *DishHeap) Push(x any) {
	var dish int
	var ok bool

	if dish, ok = x.(int); !ok {
		fmt.Fprintln(os.Stderr, "Push: type assertion failed")

		return
	}

	*h = append(*h, dish)
}

func (h *DishHeap) Pop() any {
	old := *h
	n := len(old)
	if n == 0 {
		return nil
	}
	dish := old[n-1]

	*h = old[:n-1]

	return dish
}

func main() {
	var totalDishes uint16

	if _, err := fmt.Scanln(&totalDishes); err != nil || totalDishes == 0 || totalDishes > 10000 {
		fmt.Println("invalid number of dishes")

		return
	}

	dishesHeap := &DishHeap{}
	heap.Init(dishesHeap)

	for range make([]struct{}, totalDishes) {
		var preference int

		if _, err := fmt.Scan(&preference); err != nil {
			fmt.Fprintln(os.Stderr, "couldn't read dish preference:", err)

			return
		}

		heap.Push(dishesHeap, preference)
	}

	var preferredRank uint16

	if _, err := fmt.Scan(&preferredRank); err != nil || preferredRank == 0 || dishesHeap.Len() < int(preferredRank) {
		fmt.Println("invalid preferred dish rank")

		return
	}

	var selectedDish int

	for range make([]struct{}, preferredRank) {
		value, ok := heap.Pop(dishesHeap).(int)
		if !ok {
			fmt.Fprintln(os.Stderr, "Pop: type assertion failed")

			return
		}
		selectedDish = value
	}

	fmt.Println(selectedDish)
}
