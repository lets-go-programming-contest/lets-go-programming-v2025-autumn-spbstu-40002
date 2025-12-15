package main

import (
	"container/heap"
	"errors"
	"fmt"

	intheap "github.com/manyanin.alexander/task-2-2/internal/int_heap"
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

func readDishNumber() (int, error) {
	var dishNumber int

	_, err := fmt.Scan(&dishNumber)
	if err != nil {
		return 0, ErrReadingInput
	}

	if dishNumber > MaxDishNumber || dishNumber < MinDishNumber {
		return 0, ErrIncorrectDishNumber
	}

	return dishNumber, nil
}

func readDishes(dishNumber int) ([]int, error) {
	dishes := make([]int, dishNumber)

	for dishIndex := range dishes {
		_, err := fmt.Scan(&dishes[dishIndex])
		if err != nil {
			return nil, ErrReadingInput
		}

		if dishes[dishIndex] > MaxDishRating || dishes[dishIndex] < MinDishRating {
			return nil, ErrIncorrectDishRating
		}
	}

	return dishes, nil
}

func readKValue(dishNumber int) (int, error) {
	var kValue int

	_, err := fmt.Scan(&kValue)
	if err != nil {
		return 0, ErrReadingInput
	}

	if kValue > dishNumber || kValue < 1 {
		return 0, ErrIncorrectKValue
	}

	return kValue, nil
}

func findKLargest(dishes []int, kValue int) int {
	intHeap := &intheap.IntHeap{}
	heap.Init(intHeap)

	for _, dish := range dishes {
		heap.Push(intHeap, dish)
	}

	var result int

	for range kValue {
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
	dishNumber, err := readDishNumber()
	if err != nil {
		fmt.Println(err)

		return
	}

	dishes, err := readDishes(dishNumber)
	if err != nil {
		fmt.Println(err)

		return
	}

	kValue, err := readKValue(dishNumber)
	if err != nil {
		fmt.Println(err)

		return
	}

	result := findKLargest(dishes, kValue)

	fmt.Println(result)
}
