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

func updateBounds(low, high int, op string, val int) (int, int, error) {
	if op == ">=" {
		if val > low {
			low = val
		}
	} else if op == "<=" {
		if val < high {
			high = val
		}
	} else {
		return low, high, errors.New("invalid op")
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
	for d := 0; d < n; d++ {
		var k int
		if _, err := fmt.Fscanln(os.Stdin, &k); err != nil || k < 1 || k > MaxK {
			fmt.Println(-1)
			return
		}
		low := MinTemp
		high := MaxTemp
		for i := 0; i < k; i++ {
			var op string
			var val int
			if _, err := fmt.Fscanln(os.Stdin, &op, &val); err != nil {
				fmt.Println(-1)
				return
			}
			var err error
			low, high, err = updateBounds(low, high, op, val)
			if err != nil {
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
