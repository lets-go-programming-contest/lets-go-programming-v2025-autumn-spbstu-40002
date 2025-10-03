package main

import (
	"container/heap"
	"fmt"
)

const (
	numberOfDishesMin = 1
	numberOfDishesMax = 10000
	ratingMin         = -10000
	ratingMax         = 10000
)

type MyHeap []int

func (h *MyHeap) Len() int {
	return len(*h)
}

func (h *MyHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *MyHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *MyHeap) Push(x interface{}) {
	if v, ok := x.(int); ok {
		*h = append(*h, v)
	}
}

func (h *MyHeap) Pop() interface{} {
	if len(*h) == 0 {
		return nil
	}
	
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	var numberOfDishes int

	_, err := fmt.Scan(&numberOfDishes)
	if err != nil || numberOfDishes < numberOfDishesMin || numberOfDishes > numberOfDishesMax {
		return
	}

	myHeap := &MyHeap{}
	heap.Init(myHeap)

	for range numberOfDishes {
		var rating int

		_, err = fmt.Scan(&rating)
		if err != nil || rating < ratingMin || rating > ratingMax {
			return
		}

		heap.Push(myHeap, rating)
	}

	var preferredDish int

	_, err = fmt.Scan(&preferredDish)
	if err != nil || preferredDish > numberOfDishes || preferredDish < 1 {
		return
	}

	var preferDish interface{}

	for range numberOfDishes - preferredDish + 1 {
		preferDish = heap.Pop(myHeap)
	}

	fmt.Println(preferDish)
}
