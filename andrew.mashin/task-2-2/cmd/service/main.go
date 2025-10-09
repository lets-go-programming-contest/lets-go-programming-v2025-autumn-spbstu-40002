package main

import (
	"container/heap"
	"errors"
	"fmt"

	MyHeap "github.com/Exam-Play/task-2-2/internal/min_heap"
)

const (
	numberOfDishesMin = 1
	numberOfDishesMax = 10000
	ratingMin         = -10000
	ratingMax         = 10000
)

<<<<<<< HEAD
var (
	errorIncorrectNumberOfDishes = errors.New("incorrect number of dishes")
	errorIncorrectRating         = errors.New("incorrect rating")
	errorIncorrectK              = errors.New("incorrect preferred dish")
)
=======
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
>>>>>>> 6adc7d426de11c732732c49f0c72684de4a324df

func main() {
	var numberOfDishes int

	_, err := fmt.Scan(&numberOfDishes)
	if err != nil || numberOfDishes < numberOfDishesMin || numberOfDishes > numberOfDishesMax {
		fmt.Println(errorIncorrectNumberOfDishes)

		return
	}

	myHeap := &MyHeap.MinHeap{}
	heap.Init(myHeap)

	for range numberOfDishes {
		var rating int

		_, err = fmt.Scan(&rating)
		if err != nil || rating < ratingMin || rating > ratingMax {
			fmt.Println(errorIncorrectRating)

			return
		}

		heap.Push(myHeap, rating)
	}

	var preferredDish int

	_, err = fmt.Scan(&preferredDish)
	if err != nil || preferredDish > numberOfDishes || preferredDish < 1 {
		fmt.Println(errorIncorrectK)

		return
	}

	var preferDish interface{}

	for range numberOfDishes - preferredDish + 1 {
		preferDish = heap.Pop(myHeap)
	}

	fmt.Println(preferDish)
}
