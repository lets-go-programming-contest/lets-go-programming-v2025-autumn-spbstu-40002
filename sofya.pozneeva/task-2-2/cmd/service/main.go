package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int16

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] > h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int16))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	var nDish uint16

	_, err := fmt.Scan(&nDish)
	if err != nil {
		fmt.Println("Invalid argument")
		return
	}

	rating := make(IntHeap, 0, nDish)

	heap.Init(&rating)

	for range nDish {
		var err error

		var estimation int16

		_, err = fmt.Scan(&estimation)
		if err != nil {
			fmt.Println("Invalid argument")

			return
		}

		heap.Push(&rating, estimation)
	}

	var k int16

	_, err = fmt.Scanf("\n%d\n", &k)
	if err != nil {
		fmt.Println("Invalid argument")

		return
	}

	var result int16

	for range k {
		result = heap.Pop(&rating).(int16)
	}
	fmt.Println(result)
}
