package main

import (
	"fmt"
	"os"

	"github.com/rekottt/task-2-2/kth"
	"github.com/rekottt/task-2-2/ktherr"
)

func main() {
	var itemCount int
	if _, err := fmt.Fscan(os.Stdin, &itemCount); err != nil {
		return
	}

	if itemCount < 1 || itemCount > 10000 {
		fmt.Fprintln(os.Stderr, ktherr.ErrInvalidItemCount)
		
		return
	}

	values := make([]int, itemCount)
	for i := range values {
		if _, err := fmt.Fscan(os.Stdin, &values[i]); err != nil {
			return
		}

		if values[i] < -10000 || values[i] > 10000 {
			return
		}
	}

	var position int
	if _, err := fmt.Fscan(os.Stdin, &position); err != nil {
		return
	}

	if position < 1 || position > itemCount {
		return
	}

	result, err := kth.KthMostPreferred(values, position)
	if err != nil {
		return
	}

	fmt.Println(result)
}
