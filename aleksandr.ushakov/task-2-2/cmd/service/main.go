package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/rachguta/task-2-2/myheap"
)

const (
	minNumOfDishes = 1
	maxNumOfDishes = 10000
	minMark        = -10000
	maxMark        = 10000
)

var (
	errNumOfDishes = errors.New("invalid number of dishes")
	errMark        = errors.New("invalid mark")
	errPicked      = errors.New("invalid picked dish")
)

func checkLimits(value int, minLimit int, maxLimit int) bool {
	if value >= minLimit && value <= maxLimit {
		return true
	}

	return false
}

func main() {
	var numberOfDishes int

	_, err := fmt.Scanln(&numberOfDishes)

	if err != nil || !checkLimits(numberOfDishes, minNumOfDishes, maxNumOfDishes) {
		fmt.Println(errNumOfDishes)

		return
	}

	markHeap := &myheap.IntHeap{}
	*markHeap = make(myheap.IntHeap, 0, numberOfDishes)
	heap.Init(markHeap)

	for range numberOfDishes {
		var mark int

		_, err = fmt.Scan(&mark)

		if err != nil || !checkLimits(mark, minMark, maxMark) {
			fmt.Println(errMark)

			return
		}

		heap.Push(markHeap, mark)
	}

	var picked int

	_, err = fmt.Scanln(&picked)

	if err != nil || !checkLimits(picked, 1, numberOfDishes) {
		fmt.Println(errPicked)

		return
	}

	for range picked - 1 {
		heap.Pop(markHeap)
	}

	fmt.Println(heap.Pop(markHeap))
}
