package main

import (
	"container/heap"
	"fmt"
)

type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

const (
	maxDishes = 10000
	minDishes = 1
	minScore  = -10000
	maxScore  = 10000
)

func main() {
	maxHeap := &MaxHeap{}
	heap.Init(maxHeap)

	var dishCount int

	_, err := fmt.Scan(&dishCount)
	if err != nil || dishCount < minDishes || dishCount > maxDishes {
		fmt.Println("invalid dish count")
		return
	}

	for range dishCount {
		var score int
		_, err = fmt.Scan(&score)
		if err != nil || score < minScore || score > maxScore {
			fmt.Println("invalid score value")
			return
		}

		heap.Push(maxHeap, score)
	}

	var position int
	_, err = fmt.Scan(&position)
	if err != nil {
		fmt.Println("invalid position")
		return
	}

	for range dishCount - position {
		heap.Pop(maxHeap)
	}

	result := heap.Pop(maxHeap)
	fmt.Println(result)
}
