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
	var dishNum int

	_, err := fmt.Scan(&dishNum)
	if err != nil || dishNum > maxDishNum || dishNum < minDishNum {
		fmt.Println("Error: ", errIncorrectDishNum)

		return
	}

	mainHeap := &minHeap.IntHeap{}
	heap.Init(mainHeap)

	for range dishNum {
		var curRating int
		_, err := fmt.Scan(&curRating)
		if err != nil || curRating > maxDishRating || curRating < minDishRating {
			fmt.Println("Error: ", errIncorrectDishRating)

			return
		}

		heap.Push(mainHeap, curRating)

	}

	var k int

	_, err = fmt.Scan(&k)
	if err != nil || k > dishNum || k < 1 {
		fmt.Println("Error: ", errIncorrectK)

		return
	}

	var resultRating interface{}

	for range dishNum - k + 1 {
		if mainHeap.Len() > 0 {
			resultRating = heap.Pop(mainHeap)
		}

	}

	fmt.Println(resultRating)

}
