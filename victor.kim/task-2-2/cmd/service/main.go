package main

import (
	"container/heap"
	"fmt"
	"os"
)

type IntHeap []int

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j] // максимальная куча
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x any) {
	value, ok := x.(int)
	if !ok {
		fmt.Fprintln(os.Stderr, "Push: type assertion failed")

		return
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	if n == 0 {
		return nil
	}
	x := old[n-1]

	*h = old[:n-1]

	return x
}

func main() {
	var dishCount int
	if _, err := fmt.Scan(&dishCount); err != nil {
		fmt.Fprintln(os.Stderr, "couldn't read the number of dishes:", err)

		return
	}

	dishes := &IntHeap{}
	heap.Init(dishes)

	for i := range make([]struct{}, dishCount) {
		var dish int
		if _, err := fmt.Scan(&dish); err != nil {
			fmt.Fprintln(os.Stderr, "couldn't read dish preference:", err)

			return
		}

		heap.Push(dishes, dish)
	}

	var preferredDishIndex int
	if _, err := fmt.Scan(&preferredDishIndex); err != nil || dishes.Len() < preferredDishIndex {
		fmt.Fprintln(os.Stderr, "invalid preferred dish index")

		return
	}

	var selected int
	for i := range make([]struct{}, preferredDishIndex) {
		value, ok := heap.Pop(dishes).(int)
		if !ok {
			fmt.Fprintln(os.Stderr, "Pop: type assertion failed")

			return
		}
		selected = value
	}

	fmt.Println(selected)
}
