package main

import (
	"container/heap"
	"errors"
	"fmt"

	myheap "github.com/HuaChenju/task-2-2/iternal/heap"
)

const (
	MinCount      = 1
	MaxCount      = 10000
	MinRating     = -10000
	MaxRating     = 10000
	MinPreference = 1
)

var (
	errIncorrectDishes = errors.New("invalid amount of dishes")
	errIncorrectRating = errors.New("incorrect rating")
	errIncorrectPref   = errors.New("incorrect preference")
)

func main() {
	var count, preference, rating int

	heapMax := &myheap.MaxHeap{}

	_, err := fmt.Scan(&count)

	if err != nil || count < MinCount || count > MaxCount {
		fmt.Println(errIncorrectDishes)

		return
	}

	for range count {
		_, err = fmt.Scan(&rating)

		if err != nil || rating < MinRating || rating > MaxRating {
			fmt.Println(errIncorrectRating)

			return
		}

		heap.Push(heapMax, rating)
	}

	_, err = fmt.Scan(&preference)

	if err != nil || preference > count || preference < MinPreference {
		fmt.Println(errIncorrectPref)

		return
	}

	for range preference - 1 {
		heap.Pop(heapMax)
	}

	fmt.Println(heap.Pop(heapMax))
}
