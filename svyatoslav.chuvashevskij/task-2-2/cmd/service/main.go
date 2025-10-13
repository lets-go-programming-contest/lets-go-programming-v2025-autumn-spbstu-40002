package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/Svyatoslav2324/task-2-2/internal/intheap"
)

var ErrInvalidPreferredNumber = errors.New("preferred number is more than number of dishes")

func main() {
	numberOfDishes := 0

	_, err := fmt.Scan(&numberOfDishes)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	dishes := &intheap.IntHeap{}
	heap.Init(dishes)

	var rating int
	for range numberOfDishes {
		_, err = fmt.Scan(&rating)
		if err != nil {
			fmt.Println("Error:", err)

			return
		}

		heap.Push(dishes, rating)
	}

	preferredNumber := 0

	_, err = fmt.Scan(&preferredNumber)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	if preferredNumber > numberOfDishes {
		fmt.Println("Error:", ErrInvalidPreferredNumber)

		return
	}

	for range preferredNumber - 1 {
		heap.Pop(dishes)
	}

	fmt.Println((*dishes)[0])
}
