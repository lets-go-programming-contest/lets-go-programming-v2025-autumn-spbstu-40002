package main

import (
	"container/heap"
	"fmt"
	"errors"
)

var error invalidArgument = "invalid argument"

func main() {
	var nDish uint

	_, err := fmt.Scan(&nDish)
	if err != nil {
		fmt.Println(invalidArgument)

		return
	}

	rating := make(IntHeap, 0, nDish)

	heap.Init(&rating)

	for range nDish {
		var err error

		var estimation int

		_, err = fmt.Scan(&estimation)
		if err != nil {
			fmt.Println(invalidArgument)

			return
		}

		heap.Push(&rating, estimation)
	}

	var numberOfDish int

	_, err = fmt.Scan(&numberOfDish)
	if err != nil {
		fmt.Println(invalidArgument)

		return
	}

	var result int

	for range numberOfDish {
		popped := heap.Pop(&rating)
		if value, ok := popped.(int); ok {
			result = value
		} else {
			panic(fmt.Sprintf("expected int, got %T", popped))
		}
	}

	fmt.Println(result)
}
