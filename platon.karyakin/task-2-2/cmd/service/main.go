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
		fmt.Fprintln(os.Stderr, ktherr.ErrReadItemCount)
		
		return
	}

	if itemCount < 1 || itemCount > 10000 {
		fmt.Fprintln(os.Stderr, ktherr.ErrInvalidItemCount)

		return
	}

	values := make([]int, itemCount)
	for i := range values {
		if _, err := fmt.Fscan(os.Stdin, &values[i]); err != nil {
			fmt.Fprintln(os.Stderr, ktherr.ErrReadValue)
			
			return
		}

		if values[i] < -10000 || values[i] > 10000 {
			fmt.Fprintln(os.Stderr, ktherr.ErrValueOutOfRange)
			
			return
		}
	}

	var position int
	if _, err := fmt.Fscan(os.Stdin, &position); err != nil {
		fmt.Fprintln(os.Stderr, ktherr.ErrReadPosition)
		
		return
	}

	if position < 1 || position > itemCount {
		fmt.Fprintln(os.Stderr, ktherr.ErrPositionOutOfRange)
		
		return
	}

	result, err := kth.KthMostPreferred(values, position)
	if err != nil {
		if errors.Is(err, ktherr.ErrEmptyResult) {
			fmt.Fprintln(os.Stderr, ktherr.ErrEmptyResult)
			
			return
		}
		if errors.Is(err, ktherr.ErrPositionOutOfRange) {
			fmt.Fprintln(os.Stderr, ktherr.ErrPositionOutOfRange)
			
			return
		}

		fmt.Fprintln(os.Stderr, err)
		
		return
	}

	fmt.Println(result)
}
