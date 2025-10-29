package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/ummmsh/task-2-2/intheap"
)

const (
	maxNumOfDishes = 10000
	minNumOfDishes = 1
	minRating      = -10000
	maxRating      = 10000
)

var (
	errInvalidDishes = errors.New("invalid num of dishes")
	errInvalidRating = errors.New("invalid rating value")
)

func main() {
	intHeap := &intheap.IntHeap{}
	heap.Init(intHeap)

	var numOfDishes int

	_, err := fmt.Scan(&numOfDishes)
	if err != nil || (numOfDishes < minNumOfDishes || numOfDishes > maxNumOfDishes) {
		fmt.Println(errInvalidDishes)

		return
	}

	var item int

	for range numOfDishes {
		_, err = fmt.Scan(&item)
		if err != nil || (item < minRating || item > maxRating) {
			fmt.Println(errInvalidRating)

			return
		}

		heap.Push(intHeap, item)
	}

	var rating int

	_, err = fmt.Scan(&rating)
	if err != nil {
		fmt.Println(errInvalidRating)

		return
	}

	for range numOfDishes - rating {
		heap.Pop(intHeap)
	}

	fmt.Println(heap.Pop(intHeap))
}
