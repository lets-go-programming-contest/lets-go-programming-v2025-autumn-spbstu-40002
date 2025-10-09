package main

import (
	"errors"
	"fmt"
	"os"
)

const (
	MinTemp = 15
	MaxTemp = 30
	MaxN    = 1000
	MaxK    = 1000
)

var errInvalidOp = errors.New("invalid op")

func updateBounds(low, high int, operator string, value int) (int, int, error) {
	switch operator {
	case ">=":
		if value > low {
			low = value
		}
	case "<=":
		if value < high {
			high = value
		}
	default:
		return low, high, errInvalidOp
	}

	if low < MinTemp {
		low = MinTemp
	}

	if high > MaxTemp {
		high = MaxTemp
	}

	return low, high, nil
}

func main() {
	var n int
	if _, err := fmt.Fscanln(os.Stdin, &n); err != nil || n < 1 || n > MaxN {
		fmt.Println(-1)
		return
	}

	for range make([]struct{}, n) {
		var employees int
		if _, err := fmt.Fscanln(os.Stdin, &employees); err != nil || employees < 1 || employees > MaxK {
			fmt.Println(-1)
			return
		}

		low := MinTemp

		high := MaxTemp

		for range make([]struct{}, employees) {
			var operator string
			var value int
			if _, err := fmt.Fscanln(os.Stdin, &operator, &value); err != nil {
				fmt.Println(-1)
				return
			}

			var uerr error
			low, high, uerr = updateBounds(low, high, operator, value)
			if uerr != nil {
				fmt.Println(-1)
				return
			}

			if low > high {
				fmt.Println(-1)
			} else {
				fmt.Println(low)
			}
		}
	}
}
