package utils

import (
	"container/heap"

	"github.com/slendycs/task-2-2/interfaces"
)

func GetPreferredDish(pickedDishId int, dishesRaitingList []int) int {
	// Initialize the heap.
	raitingHeap := &interfaces.MinHeap{}
	heap.Init(raitingHeap)

	// Add each rating to the heap, pushing the minimum one out of it if its size exceeds the ID of the selected dish.
	for _, raiting := range dishesRaitingList {
		heap.Push(raitingHeap, raiting)
		if raitingHeap.Len() > pickedDishId {
			_ = heap.Pop(raitingHeap)
		}
	}

	return (*raitingHeap)[0]
}
