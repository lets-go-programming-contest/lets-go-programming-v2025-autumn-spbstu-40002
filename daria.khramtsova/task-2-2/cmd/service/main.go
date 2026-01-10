package main

import (
	"container/heap"
	"errors"
	"fmt"

	maxheap "github.com/hehemka/task-2-2/internal/heap"
)

const (
	minNum             = 1
	maxNum             = 10000
	minRating          = -10000
	maxRating          = 10000
	minPreferredNumber = 1
)

var (
	errIncorrectDishNumber = errors.New("incorrect number of dishes")
	errIncorrectRating     = errors.New("incorrect rating")
	errIncorrectPreference = errors.New("incorrect preferred number")
)

func main() {
	maxHeap := &maxheap.MaxHeap{}
	heap.Init(maxHeap)

	var numberOfDishes, preferredNum, dishRating int

	_, err := fmt.Scan(&numberOfDishes)
	if err != nil || (numberOfDishes < minNum || numberOfDishes > maxNum) {
		fmt.Println(errIncorrectDishNumber)

		return
	}

	for range numberOfDishes {
		_, err = fmt.Scan(&dishRating)
		if err != nil || (dishRating < minRating || dishRating > maxRating) {
			fmt.Println(errIncorrectRating)

			return
		}

		heap.Push(maxHeap, dishRating)
	}

	_, err = fmt.Scan(&preferredNum)
	if err != nil || preferredNum > numberOfDishes || preferredNum < minPreferredNumber {
		fmt.Println(errIncorrectPreference)

		return
	}

	for range preferredNum - 1 {
		heap.Pop(maxHeap)
	}

	fmt.Println(heap.Pop(maxHeap))
}
