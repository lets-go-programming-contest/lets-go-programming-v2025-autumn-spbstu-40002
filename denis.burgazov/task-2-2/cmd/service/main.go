package main

import (
	"container/heap"
	"fmt"
)

func parse(a any, errText string) bool {
	_, err := fmt.Scanln(a)
	if err != nil {
		fmt.Println(errText)

		var errBuf string

		_, _ = fmt.Scanln(&errBuf)
	}

	return err == nil
}

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	n := len(*h)
	x := (*h)[n-1]
	*h = (*h)[:n-1]

	return x
}

func findKthLargest(nums []int, k int) int {
	h := &IntHeap{}
	heap.Init(h)

	for i := range k {
		heap.Push(h, nums[i])
	}

	for i := k; i < len(nums); i++ {
		if nums[i] > (*h)[0] {
			heap.Pop(h)
			heap.Push(h, nums[i])
		}
	}

	return (*h)[0]
}

func main() {
	var numberOfDishes, k int

	if !parse(&numberOfDishes, "Invalid number of dishes") {
		return
	}

	ratings := make([]int, numberOfDishes)
	for i := range numberOfDishes {
		fmt.Scan(&ratings[i])
	}

	if !parse(&k, "Invalid k-number") {
		return
	}

	fmt.Println(findKthLargest(ratings, k))
}
