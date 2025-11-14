package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/sonychello/task-2-2/intheap"
)

var errInvalidArgument = errors.New("invalid argument")

func main() {
	var nDish uint16

	_, err := fmt.Scan(&nDish)
	if err != nil {
		fmt.Println(errInvalidArgument)

		return
	}

	rating := make(intheap.IntHeap, 0, nDish)

	heap.Init(&rating)

	for range nDish {
		var err error

		var estimation int16

		_, err = fmt.Scan(&estimation)
		if err != nil {
			fmt.Println(errInvalidArgument)

			return
		}

		heap.Push(&rating, estimation)
	}

	var numberOfDish int16

	_, err = fmt.Scan(&numberOfDish)
	if err != nil {
		fmt.Println(errInvalidArgument)

		return
	}

	var result int16

	for range numberOfDish {
		popped := heap.Pop(&rating)
		if value, ok := popped.(int16); ok {
			result = value
		} else {
			panic(fmt.Sprintf("expected int16, got %T", popped))
		}
	}

	fmt.Println(result)
}
