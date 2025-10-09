package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[j] < (*h)[i] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x any) {
	v, ok := x.(int)
	if !ok {
		panic("IntHeap can only contain ints")
	}
	*h = append(*h, v)
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]

	return x
}

func main() {
	numberOfDishes := 0

	_, err := fmt.Scan(&numberOfDishes)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	dishes := &IntHeap{}
	heap.Init(dishes)

	var rating int
	for range numberOfDishes {
		_, err = fmt.Scan(&rating)
		if err != nil {
			fmt.Println("Error:", err)

			return
		}

		heap.Push(dishes, rating)
	}

	preferredNumber := 0

	_, err = fmt.Scan(&preferredNumber)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	if preferredNumber > numberOfDishes {
		fmt.Println("Error: preferred number is more than number of dishes")

		return
	}

	for range preferredNumber - 1 {
		heap.Pop(dishes)
	}

	fmt.Println((*dishes)[0])
}
