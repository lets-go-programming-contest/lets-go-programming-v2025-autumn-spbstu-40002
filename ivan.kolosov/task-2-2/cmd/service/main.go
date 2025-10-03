package main

import (
	"container/heap"
	"errors"
	"fmt"
)

const (
	AmountOfDishesMax  = 10000
	AmountOfDishesMin  = 1
	RatingOfDishesMax  = 10000
	RatingOfDishesMin  = -10000
	NumberOfTheDishMin = 1
)

var (
	errIncorrectAmountOfDishes   = errors.New("incorrect amount of dishes")
	errIncorrectRatingForTheDish = errors.New("incorrect rating for the dish")
	errIncorrectK                = errors.New("incorrect k")
)

type MaxHeap []int

func (heap *MaxHeap) Len() int { return len((*heap)) }

func (heap *MaxHeap) Less(i, j int) bool { return (*heap)[i] > (*heap)[j] }

func (heap *MaxHeap) Swap(i, j int) { (*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i] }

func (heap *MaxHeap) Push(x interface{}) {
	value, isInt := x.(int)
	if !isInt {
		panic("expected int")
	}

	*heap = append(*heap, value)
}

func (heap *MaxHeap) Pop() interface{} {
	if len(*heap) == 0 {
		panic("trying to pop from an empty heap")
	}

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

	heapOfDishes := &MaxHeap{}
	heap.Init(heapOfDishes)

	for range amountOfDishes {
		_, err = fmt.Scan(&newDish)
		if err != nil || newDish > RatingOfDishesMax || newDish < RatingOfDishesMin {
			fmt.Println("Error:", errIncorrectRatingForTheDish)

			return
		}

		heap.Push(heapOfDishes, newDish)
	}

	var numberOfTheDish int

	_, err = fmt.Scan(&numberOfTheDish)
	if err != nil || numberOfTheDish > amountOfDishes || numberOfTheDish < NumberOfTheDishMin {
		fmt.Println("Error:", errIncorrectK)

		return
	}

	var todaysDish interface{}

	for range numberOfTheDish {
		todaysDish = heap.Pop(heapOfDishes)
	}

	fmt.Println(todaysDish)
}
