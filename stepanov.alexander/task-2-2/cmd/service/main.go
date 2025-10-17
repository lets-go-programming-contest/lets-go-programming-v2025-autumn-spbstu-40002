package main

import (
	"container/heap"
	"fmt"
)

// MaxHeap представляет максимальную кучу для целых чисел
type MaxHeap []int

func (h MaxHeap) Len() int {
	return len(h)
}

func (h MaxHeap) Less(i, j int) bool {
	return h[i] > h[j]
}

func (h MaxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	var N int
	fmt.Scan(&N)

	dishes := make([]int, N)
	for i := 0; i < N; i++ {
		fmt.Scan(&dishes[i])
	}

	var k int
	fmt.Scan(&k)

	h := &MaxHeap{}
	heap.Init(h)

	for _, dish := range dishes {
		heap.Push(h, dish)
	}

	var result int
	for i := 0; i < k; i++ {
		result = heap.Pop(h).(int)
	}

	fmt.Println(result)
}
