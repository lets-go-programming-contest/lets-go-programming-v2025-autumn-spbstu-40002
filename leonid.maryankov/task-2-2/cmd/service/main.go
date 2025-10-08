package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x any) {
	val, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, val)
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

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
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	dishHeap := IntHeap(dishes)
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
