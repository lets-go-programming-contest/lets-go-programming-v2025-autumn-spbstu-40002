package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/rekottt/task-2-2/kth"
	"github.com/rekottt/task-2-2/ktherr"
)

func readItemCount() (int, error) {
	var itemCount int
	if _, err := fmt.Fscan(os.Stdin, &itemCount); err != nil {
		return 0, ktherr.ErrReadItemCount
	}

	if itemCount < 1 || itemCount > 10000 {
		return 0, ktherr.ErrInvalidItemCount
	}

	return itemCount, nil
}

func readValues(count int) ([]int, error) {
	values := make([]int, count)
	for i := range values {
		if _, err := fmt.Fscan(os.Stdin, &values[i]); err != nil {
			return nil, ktherr.ErrReadValue
		}

		if values[i] < -10000 || values[i] > 10000 {
			return nil, ktherr.ErrValueOutOfRange
		}
	}

	return values, nil
}

func readPosition(itemCount int) (int, error) {
	var position int
	if _, err := fmt.Fscan(os.Stdin, &position); err != nil {
		return 0, ktherr.ErrReadPosition
	}

	if position < 1 || position > itemCount {
		return 0, ktherr.ErrPositionOutOfRange
	}

	return position, nil
}

func main() {
	itemCount, err := readItemCount()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		return
	}

	values, err := readValues(itemCount)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		return
	}

	position, err := readPosition(itemCount)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		return
	}

	result, err := kth.KthMostPreferred(values, position)
	if err != nil {
		switch {
		case errors.Is(err, ktherr.ErrEmptyResult):
			fmt.Fprintln(os.Stderr, ktherr.ErrEmptyResult)

			return
		case errors.Is(err, ktherr.ErrPositionOutOfRange):
			fmt.Fprintln(os.Stderr, ktherr.ErrPositionOutOfRange)

			return
		default:
			fmt.Fprintln(os.Stderr, err)

			return
		}
	}

	fmt.Println(result)
}
