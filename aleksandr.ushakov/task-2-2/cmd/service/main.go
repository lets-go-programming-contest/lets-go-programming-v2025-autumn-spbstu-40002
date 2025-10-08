package main

import (
	"container/heap"
	"fmt"
)

const (
	minNumOfDishes = 1
	maxNumOfDishes = 10000
	minMark        = -10000
	maxMark        = 10000
)

// Heap realization
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

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
		return
	}
	markHeap := &IntHeap{}
	*markHeap = make(IntHeap, 0, numberOfDishes)
	heap.Init(markHeap)
	for range numberOfDishes {
		var mark int
		_, err = fmt.Scan(&mark)
		if err != nil || !checkLimits(mark, minMark, maxMark) {
			return
		}
		heap.Push(markHeap, mark)
	}
	var k int
	_, err = fmt.Scanln(&k)
	if err != nil || !checkLimits(k, 1, numberOfDishes) {
		return
	}
	for range k - 1 {
		heap.Pop(markHeap)
	}
	fmt.Println(heap.Pop(markHeap))

}
