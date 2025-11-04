package main

import (
	"container/heap"
	"fmt"
)

type PriorityQueue []int

func (pq *PriorityQueue) Len() int { return len(*pq) }

func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i] > (*pq)[j]
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		return
	}

	*pq = append(*pq, value)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]

	return item
}

const (
	maxDishesCount = 10000
	minDishesCount = 1
	minValue       = -10000
	maxValue       = 10000
)

func main() {
	priorityQueue := &PriorityQueue{}
	heap.Init(priorityQueue)

	var totalDishes int

	_, err := fmt.Scan(&totalDishes)
	if err != nil {
		fmt.Println("invalid input")

		return
	}

	if totalDishes < minDishesCount || totalDishes > maxDishesCount {
		fmt.Println("invalid dish count")

		return
	}

	for range totalDishes {
		var currentValue int

		_, err = fmt.Scan(&currentValue)
		if err != nil {
			fmt.Println("invalid input")

			return
		}

		if currentValue < minValue || currentValue > maxValue {
			fmt.Println("invalid value")

			return
		}

		heap.Push(priorityQueue, currentValue)
	}

	var targetPosition int

	_, err = fmt.Scan(&targetPosition)
	if err != nil {
		fmt.Println("invalid input")

		return
	}

	for range targetPosition - 1 {
		heap.Pop(priorityQueue)
	}

	result := heap.Pop(priorityQueue)
	fmt.Println(result)
}
