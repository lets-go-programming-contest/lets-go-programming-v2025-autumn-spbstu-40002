package main

import (
	"container/heap"
	"errors"
	"fmt"

	minHeap "github.com/ControlShiftEscape/task-2-2/internal/minheap"
)

const (
	minDishNum    = 1
	maxDishNum    = 10000
	minDishRating = -10000
	maxDishRating = 10000
)

var (
	errIncorrectDishNum    = errors.New("incorrect dish number")
	errIncorrectDishRating = errors.New("incorrect dish rating")
	errIncorrectK          = errors.New("incorrect value for k")
)

func main() {
	dishNum, err := readDishNumber()
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	mainHeap, err := readAndCreateHeap(dishNum)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	kValue, err := readKValue(dishNum)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	result := clearHeapUntilK(mainHeap, dishNum, kValue)
	fmt.Println(result)
}

func readDishNumber() (int, error) {
	var dishNum int

	_, err := fmt.Scan(&dishNum)
	if err != nil || dishNum > maxDishNum || dishNum < minDishNum {
		return 0, errIncorrectDishNum
	}

	return dishNum, nil
}

func readAndCreateHeap(dishNum int) (*minHeap.IntHeap, error) {
	mainHeap := &minHeap.IntHeap{}
	heap.Init(mainHeap)

	for range dishNum {
		var curRating int
		_, err := fmt.Scan(&curRating)

		if err != nil || curRating > maxDishRating || curRating < minDishRating {
			return nil, errIncorrectDishRating
		}

		heap.Push(mainHeap, curRating)
	}

	return mainHeap, nil
}

func readKValue(dishNum int) (int, error) {
	var kValue int
	_, err := fmt.Scan(&kValue)

	if err != nil || kValue > dishNum || kValue < 1 {
		return 0, errIncorrectK
	}

	return kValue, nil
}

func clearHeapUntilK(mainHeap *minHeap.IntHeap, dishNum, kValue int) interface{} {
	var result interface{}

	for range dishNum - kValue + 1 {
		if mainHeap.Len() > 0 {
			result = heap.Pop(mainHeap)
		}
	}

	return result
}
