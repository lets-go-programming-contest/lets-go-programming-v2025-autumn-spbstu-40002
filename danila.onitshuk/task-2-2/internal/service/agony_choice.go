package service

import (
	"container/heap"
	"fmt"

	inputError "danila.onitshuk/task-2-2/internal/errors"
	interfaces "danila.onitshuk/task-2-2/internal/heap"
)

func AgonyChoice() {
	dishes := &interfaces.Heap{}
	heap.Init(dishes)

	var (
		dishCount        int
		dishPreferenceID int
	)

	_, err := fmt.Scan(&dishCount)
	if err != nil {
		fmt.Println(inputError.ErrInvalidTypeData)

		return
	} else if dishCount < 1 || dishCount > 10_000 {
		fmt.Println(inputError.ErrInvalidNumberDishes)

		return
	}

	for i := 0; i < dishCount; i++ {
		var dish int

		_, err := fmt.Scan(&dish)
		if err != nil {
			fmt.Println(inputError.ErrInvalidTypeData)

			return
		} else if dish < -10_000 || dish > 10_000 {
			fmt.Println(inputError.ErrInvalidPriorityDish)

			return
		}

		heap.Push(dishes, dish)
	}

	_, err = fmt.Scan(&dishPreferenceID)
	if err != nil {
		fmt.Println(inputError.ErrInvalidTypeData)

		return
	} else if dishPreferenceID < 1 || dishPreferenceID > dishCount {
		fmt.Println(inputError.ErrInvalidNumberDishes)

		return
	}

	for range dishPreferenceID - 1 {
		heap.Pop(dishes)
	}

	fmt.Println((*dishes)[0])
}
