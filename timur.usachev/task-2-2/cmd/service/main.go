package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/t1wt/task-2-2/internal/interheap"
)

var (
	errInvalidCount     = errors.New("invalid input: count must be a positive int")
	errFailedReadValues = errors.New("invalid input: failed to read values :(")
	errInvalidKth       = errors.New("invalid input: 1 <= k <= count")
)

func readInput() ([]int, int) {
	var count int
	if _, err := fmt.Fscan(os.Stdin, &count); err != nil || count <= 0 {
		fmt.Fprintln(os.Stderr, errInvalidCount)
		os.Exit(1)
	}

	values := make([]int, count)
	for i := range values {
		if _, err := fmt.Fscan(os.Stdin, &values[i]); err != nil {
			fmt.Fprintln(os.Stderr, errFailedReadValues)
			os.Exit(1)
		}
	}

	var kth int
	if _, err := fmt.Fscan(os.Stdin, &kth); err != nil || kth <= 0 || kth > count {
		fmt.Fprintln(os.Stderr, errInvalidKth)
		os.Exit(1)
	}

	return values, kth
}

func main() {
	values, kth := readInput()

	result := interheap.FindKthLargest(values, kth)
	fmt.Println(result)
}
