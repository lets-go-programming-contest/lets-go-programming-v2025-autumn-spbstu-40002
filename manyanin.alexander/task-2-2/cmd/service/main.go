package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/manyanin.alexander/task-2-2/internal/int_heap"
)

const (
	MinDishNumber = 1
	MaxDishNumber = 10000
	MinDishRating = -10000
	MaxDishRating = 10000
)

var (
	ErrIncorrectDishNumber = errors.New("incorrect dish number")
	ErrIncorrectDishRating = errors.New("incorrect dish rating")
	ErrIncorrectKValue     = errors.New("incorrect k-value")
	ErrReadingInput        = errors.New("error reading input")
	ErrUnexpectedType      = errors.New("unexpected type")
)

func readDishNumber() int {
	var dishNumber int

	_, err := fmt.Scan(&dishNumber)
	if err != nil {
		panic(ErrReadingInput)
	}

	if dishNumber > MaxDishNumber || dishNumber < MinDishNumber {
		panic(ErrIncorrectDishNumber)
	}

	return dishNumber
}

func readDishes(dishNumber int) []int {
	dishes := make([]int, dishNumber)

	for dishIndex := 0; dishIndex < dishNumber; dishIndex++ {
		_, err := fmt.Scan(&dishes[dishIndex])
		if err != nil {
			panic(ErrReadingInput)
		}

		if dishes[dishIndex] > MaxDishRating || dishes[dishIndex] < MinDishRating {
			panic(ErrIncorrectDishRating)
		}
	}

	return dishes
}

func readKValue(dishNumber int) int {
	var kValue int

	_, err := fmt.Scan(&kValue)
	if err != nil {
		panic(ErrReadingInput)
	}

	if kValue > dishNumber || kValue < 1 {
		panic(ErrIncorrectKValue)
	}

	return kValue
}

func findKLargest(dishes []int, kValue int) int {
	intHeap := &int_heap.IntHeap{}

	heap.Init(intHeap)

	for _, dish := range dishes {
		heap.Push(intHeap, dish)
	}

	var result int

	for index := 0; index < kValue; index++ {
		poppedValue := heap.Pop(intHeap)

		value, correctType := poppedValue.(int)

		if !correctType {
			panic(ErrUnexpectedType)
		}

		result = value
	}

	return result
}

func main() {
	dishNumber := readDishNumber()

	dishes := readDishes(dishNumber)

	kValue := readKValue(dishNumber)

	result := findKLargest(dishes, kValue)

	fmt.Println(result)
}
