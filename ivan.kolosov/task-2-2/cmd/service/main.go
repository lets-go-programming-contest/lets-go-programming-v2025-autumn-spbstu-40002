package main

import (
	"container/heap"
	"errors"
	"fmt"

	myHeap "github.com/InsomniaDemon/task-2-2/internal/max_heap"
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

func main() {
	var amountOfDishes int

	_, err := fmt.Scan(&amountOfDishes)
	if err != nil || amountOfDishes > AmountOfDishesMax || amountOfDishes < AmountOfDishesMin {
		fmt.Println("Error:", errIncorrectAmountOfDishes)

		return
	}

	var newDish int

	heapOfDishes := &myHeap.MaxHeap{}
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
