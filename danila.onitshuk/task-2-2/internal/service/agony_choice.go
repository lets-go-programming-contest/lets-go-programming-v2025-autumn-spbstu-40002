package service

import (
	"container/heap"
	"fmt"

	inputerror "danila.onitshuk/task-2-2/internal/errors"
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
		fmt.Println(inputerror.ErrInvalidTypeData)

		return
	} else if dishCount < 1 || dishCount > 10_000 {
		fmt.Println(inputerror.ErrInvalidNumberDishes)

		return
	}

	for range dishCount {
		var dish int

		_, err := fmt.Scan(&dish)
		if err != nil {
			fmt.Println(inputerror.ErrInvalidTypeData)

			return
		} else if dish < -10_000 || dish > 10_000 {
			fmt.Println(inputerror.ErrInvalidPriorityDish)

			return
		}

		heap.Push(dishes, dish)
	}

	_, err = fmt.Scan(&dishPreferenceID)
	if err != nil {
		fmt.Println(inputerror.ErrInvalidTypeData)

		return
	} else if dishPreferenceID < 1 || dishPreferenceID > dishCount {
		fmt.Println(inputerror.ErrInvalidNumberDishes)

		return
	}

	for range dishPreferenceID - 1 {
		heap.Pop(dishes)
	}

	fmt.Println((*dishes)[0])
}
