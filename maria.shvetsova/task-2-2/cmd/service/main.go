package main

import (
	"container/heap"
	"fmt"

	"github.com/ummmsh/task-2-2/intheap"
)

const (
	maxNumOfDishes = 10000
	minNumOfDishes = 1
	minRating      = -10000
	maxRating      = 10000
)

func main() {
	intHeap := &intheap.IntHeap{}
	heap.Init(intHeap)

	var numOfDishes int

	_, err := fmt.Scan(&numOfDishes)
	if err != nil {
		return
	}

	if numOfDishes < minNumOfDishes || numOfDishes > maxNumOfDishes {
		return
	}

	var item int

	for range numOfDishes {
		_, err = fmt.Scan(&item)
		if err != nil {
			return
		}

		if item < minRating || item > maxRating {
			return
		}

		heap.Push(intHeap, item)
	}

	var rating int

	_, err = fmt.Scan(&rating)
	if err != nil {
		return
	}

	for range numOfDishes - rating {
		heap.Pop(intHeap)
	}

	fmt.Println(heap.Pop(intHeap))
}
