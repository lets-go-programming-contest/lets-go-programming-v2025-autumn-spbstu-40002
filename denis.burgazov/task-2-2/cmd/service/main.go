package main

import (
	"container/heap"
	"fmt"
)

func parse(a any, errText string) bool {
	_, err := fmt.Scan(a)
	if err != nil {
		fmt.Println(errText)

		var errBuf string

		_, _ = fmt.Scanln(&errBuf)
	}

	return err == nil
}

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x interface{}) {
	if num, ok := x.(int); ok {
		*h = append(*h, num)
	}
}

func (h *IntHeap) Pop() interface{} {
	n := len(*h)
	x := (*h)[n-1]
	*h = (*h)[:n-1]

	return x
}

func findKthLargest(nums []int, targetIndex int) int {
	minHeap := &IntHeap{}
	heap.Init(minHeap)

	for i := range targetIndex {
		heap.Push(minHeap, nums[i])
	}

	for i := targetIndex; i < len(nums); i++ {
		if nums[i] > (*minHeap)[0] {
			heap.Pop(minHeap)
			heap.Push(minHeap, nums[i])
		}
	}

	return (*minHeap)[0]
}

func main() {
	var numberOfDishes, targetIndex int

	if !parse(&numberOfDishes, "Invalid number of dishes") {
		return
	}

	ratings := make([]int, numberOfDishes)
	for i := range numberOfDishes {
		if !parse(&ratings[i], "Invalid rating") {
			return
		}
	}

	if !parse(&targetIndex, "Invalid k-number") {
		return
	}

	fmt.Println(findKthLargest(ratings, targetIndex))
}
