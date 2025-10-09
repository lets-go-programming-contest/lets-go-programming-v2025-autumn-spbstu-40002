package main

import (
	"container/heap"
	"errors"
	"fmt"

	MyHeap "github.com/Exam-Play/task-2-2/internal/min_heap"
)

const (
	numberOfDishesMin = 1
	numberOfDishesMax = 10000
	ratingMin         = -10000
	ratingMax         = 10000
)

var (
	errIncorrectNumberOfDishes = errors.New("incorrect number of dishes")
	errIncorrectRating         = errors.New("incorrect rating")
	errIncorrectK              = errors.New("incorrect preferred dish")
)

func main() {
	var numberOfDishes int

	_, err := fmt.Scan(&numberOfDishes)
	if err != nil || numberOfDishes < numberOfDishesMin || numberOfDishes > numberOfDishesMax {
		fmt.Println(errIncorrectNumberOfDishes)

		return
	}

	myHeap := &MyHeap.MinHeap{}
	heap.Init(myHeap)

	for range numberOfDishes {
		var rating int

		_, err = fmt.Scan(&rating)
		if err != nil || rating < ratingMin || rating > ratingMax {
			fmt.Println(errIncorrectRating)

			return
		}

		heap.Push(myHeap, rating)
	}

	var preferredDish int

	_, err = fmt.Scan(&preferredDish)
	if err != nil || preferredDish > numberOfDishes || preferredDish < 1 {
		fmt.Println(errIncorrectK)

		return
	}

	var preferDish interface{}

	for range numberOfDishes - preferredDish + 1 {
		preferDish = heap.Pop(myHeap)
	}

	fmt.Println(preferDish)
}
