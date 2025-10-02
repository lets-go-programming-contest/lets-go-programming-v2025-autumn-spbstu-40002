package main

import (
	"container/heap"
	"errors"
	"fmt"
)

const (
	AmountOfDishesMax = 10000
	AmountOfDishesMin = 1
	RatingOfDishesMax = 10000
	RatingOfDishesMin = -10000
)

var (
	errIncorrectAmountOfDishes   = errors.New("incorrect amount of dishes")
	errIncorrectRatingForTheDish = errors.New("incorrect rating for the dish")
	errIncorrectK                = errors.New("incorrect k")
)

type MaxHeap []int

func (heap MaxHeap) Len() int { return len(heap) }

func (heap MaxHeap) Less(i, j int) bool { return heap[i] > heap[j] }

func (heap MaxHeap) Swap(i, j int) { heap[i], heap[j] = heap[j], heap[i] }

func (heap *MaxHeap) Push(x interface{}) { *heap = append(*heap, x.(int)) }

func (heap *MaxHeap) Pop() interface{} {
	x := (*heap)[len(*heap)-1]
	*heap = (*heap)[0 : len(*heap)-1]
	return x
}

func main() {
	var amountOfDishes int

	_, err := fmt.Scan(&amountOfDishes)
	if err != nil || amountOfDishes > AmountOfDishesMax || amountOfDishes < AmountOfDishesMin {
		fmt.Println("Error:", errIncorrectAmountOfDishes)

		return
	}

	var newDish int

	myHeap := &MaxHeap{}
	heap.Init(myHeap)

	for range amountOfDishes {
		_, err = fmt.Scan(&newDish)
		if err != nil || newDish > RatingOfDishesMax || newDish < RatingOfDishesMin {
			fmt.Println("Error:", errIncorrectRatingForTheDish)

			return
		}
		heap.Push(myHeap, newDish)
	}

	var k int

	_, err = fmt.Scan(&k)
	if err != nil || k > amountOfDishes {
		fmt.Println("Error:", errIncorrectK)

		return
	}

	fmt.Println((*myHeap)[k-1])
}
