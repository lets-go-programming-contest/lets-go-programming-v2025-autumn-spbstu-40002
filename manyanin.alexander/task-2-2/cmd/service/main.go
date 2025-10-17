package main

import (
	"container/heap"
	"errors"
	"fmt"
)

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] > h[j] // max-heap получается; если < - min-heap
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
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

func main() {
	var dishNum int
	_, err := fmt.Scan(&dishNum)
	if err != nil || dishNum > MaxDishNum || dishNum < MinDishNum {
		fmt.Printf("Error reading input: %v\n", ErrIncorrectDishNum)
		return
	}

	dishes := make([]int, dishNum)
	for i := 0; i < dishNum; i++ {
		_, err := fmt.Scan(&dishes[i])
		if err != nil || dishes[i] > MaxDishRating || dishes[i] < MinDishRating {
			fmt.Printf("Error reading input: %v\n", ErrIncorrectDishRate)
			return
		}
	}

	var kValue int
	_, err = fmt.Scan(&kValue)
	if err != nil || kValue > dishNum || kValue < 1 {
		fmt.Printf("Error reading input: %v\n", ErrIncorrectK)
		return
	}

	h := &IntHeap{}
	heap.Init(h)

	for _, dish := range dishes {
		heap.Push(h, dish)
	}

	var result int
	for i := 0; i < kValue; i++ {
		result = heap.Pop(h).(int)
	}

	fmt.Println(result)
}
