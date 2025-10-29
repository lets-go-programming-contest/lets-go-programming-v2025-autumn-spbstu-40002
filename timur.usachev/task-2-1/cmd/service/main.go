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

type TemperatureRange struct {
	low  int
	high int
}

func (t *TemperatureRange) update(operator string, value int) error {
	switch operator {
	case ">=":
		if value > t.low {
			t.low = value
		}
	case "<=":
		if value < t.high {
			t.high = value
		}
	default:
		return errInvalidOp
	}

	if t.low < MinTemp {
		t.low = MinTemp
	}

	if t.high > MaxTemp {
		t.high = MaxTemp
	}

	return nil
}

func (t *TemperatureRange) isInvalid() bool {
	return t.low > t.high
}

func main() {
	var numCases int

	if _, err := fmt.Fscanln(os.Stdin, &numCases); err != nil || numCases < 1 || numCases > MaxN {
		fmt.Println("invalid number of cases")

		return
	}

	for range make([]struct{}, numCases) {
		var employees int

		if _, err := fmt.Fscanln(os.Stdin, &employees); err != nil || employees < 1 || employees > MaxK {
			fmt.Println("invalid number of employees")

			return
		}

		tempRange := TemperatureRange{low: MinTemp, high: MaxTemp}

		for range make([]struct{}, employees) {
			var (
				operator string
				value    int
			)

			if _, err := fmt.Fscanln(os.Stdin, &operator, &value); err != nil {
				fmt.Println("input error")

				return
			}

			if err := tempRange.update(operator, value); err != nil {
				fmt.Println("invalid operation")

				return
			}

			if tempRange.isInvalid() {
				fmt.Println(-1)
			} else {
				fmt.Println(tempRange.low)
			}
		}
	}
}
