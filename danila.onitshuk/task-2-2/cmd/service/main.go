package main

import (
	"container/heap"
	"fmt"

	"danila.onitshuk/task-2-2/internal/buffet"
)

func main() {
	dishes := &buffet.Heap{}
	heap.Init(dishes)

	var (
		dishCount        int
		dishPreferenceID int
	)

	_, err := fmt.Scan(&dishCount)
	if err != nil {
		fmt.Println(buffet.ErrInvalidTypeData)

		return
	} else if dishCount < 1 || dishCount > 10_000 {
		fmt.Println(buffet.ErrInvalidNumberDishes)

		return
	}

	for range dishCount {
		var dish int

		_, err := fmt.Scan(&dish)
		if err != nil {
			fmt.Println(buffet.ErrInvalidTypeData)

			return
		} else if dish < -10_000 || dish > 10_000 {
			fmt.Println(buffet.ErrInvalidPriorityDish)

			return
		}

		heap.Push(dishes, dish)
	}

	_, err = fmt.Scan(&dishPreferenceID)
	if err != nil {
		fmt.Println(buffet.ErrInvalidTypeData)

		return
	}

	buffet.ChoosingDish(dishPreferenceID, dishes)
}
