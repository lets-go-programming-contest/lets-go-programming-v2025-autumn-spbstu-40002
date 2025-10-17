package main

import (
	"container/heap"
	"errors"
	"fmt"
)

type IntHeap []int

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j] // max-heap получается; если < - min-heap
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

const (
	MinDishNum    = 1
	MaxDishNum    = 10000
	MinDishRating = -10000
	MaxDishRating = 10000
)

var (
	ErrIncorrectDishNum  = errors.New("incorrect dish number")
	ErrIncorrectDishRate = errors.New("incorrect dish rating")
	ErrIncorrectK        = errors.New("incorrect k-value")
)

func readDishNumber() (int, error) {
	var dishNum int

	_, err := fmt.Scan(&dishNum)
	if err != nil {
		return 0, fmt.Errorf("error reading input: %w", err)
	}

	if dishNum > MaxDishNum || dishNum < MinDishNum {
		return 0, ErrIncorrectDishNum
	}

	return dishNum, nil
}

func readDishes(dishNum int) ([]int, error) {
	dishes := make([]int, dishNum)

	for i := range dishNum {
		_, err := fmt.Scan(&dishes[i])
		if err != nil {
			return nil, fmt.Errorf("error reading input: %w", err)
		}

		if dishes[i] > MaxDishRating || dishes[i] < MinDishRating {
			return nil, ErrIncorrectDishRate
		}
	}

	return dishes, nil
}

func readKValue(dishNum int) (int, error) {
	var kValue int

	_, err := fmt.Scan(&kValue)
	if err != nil {
		return 0, fmt.Errorf("error reading input: %w", err)
	}

	if kValue > dishNum || kValue < 1 {
		return 0, ErrIncorrectK
	}

	return kValue, nil
}

func findKLargest(dishes []int, kValue int) int {
	intHeap := &IntHeap{}
	heap.Init(intHeap)

	for _, dish := range dishes {
		heap.Push(intHeap, dish)
	}

	var result int

	for range kValue {
		poppedValue := heap.Pop(intHeap)
		value, ok := poppedValue.(int)

		if !ok {
			continue
		}

		result = value
	}

	return result
}

func main() {
	dishNum, err := readDishNumber()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	dishes, err := readDishes(dishNum)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	kValue, err := readKValue(dishNum)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	result := findKLargest(dishes, kValue)
	fmt.Println(result)
}
