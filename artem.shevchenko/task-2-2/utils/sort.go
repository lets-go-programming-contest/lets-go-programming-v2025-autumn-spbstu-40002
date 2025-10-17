package utils

import (
	"container/heap"

	"github.com/slendycs/task-2-2/interfaces"
)

func GetPreferredDish(pickedDishID int, dishesRatingList []int) int {
	// Initialize the heap.
	ratingHeap := &interfaces.MinHeap{}
	heap.Init(ratingHeap)

	// Add each rating to the heap, pushing the minimum one out of it if its size exceeds the ID of the selected dish.
	for _, rating := range dishesRatingList {
		heap.Push(ratingHeap, rating)

		if ratingHeap.Len() > pickedDishID {
			_ = heap.Pop(ratingHeap)
		}
	}

	return (*ratingHeap)[0]
}
