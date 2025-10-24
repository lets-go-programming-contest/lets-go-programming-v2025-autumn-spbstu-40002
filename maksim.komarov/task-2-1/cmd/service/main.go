package main

import (
	"errors"
	"fmt"
)

const (
	minTemp = 15
	maxTemp = 30
)

var (
	ErrUnknownOperator = errors.New("unknown operator")
	ErrNoFeasibleTemp  = errors.New("no feasible temperature")
)

type TempManager struct {
	low  int
	high int
}

func NewTempManager() *TempManager {
	return &TempManager{low: minTemp, high: maxTemp}
}

func (t *TempManager) Apply(operatorSign string, value int) error {
	switch operatorSign {
	case ">=":
		if value > t.low {
			t.low = value
		}
	case "<=":
		if value < t.high {
			t.high = value
		}
	default:
		return fmt.Errorf("%w: %q", ErrUnknownOperator, operatorSign)
	}

	return nil
}

func (t *TempManager) Current() (int, error) {
	if t.low > t.high {
		return -1, ErrNoFeasibleTemp
	}

	return t.low, nil
}

func main() {
	var departmentCount int
	if _, err := fmt.Scan(&departmentCount); err != nil {
		fmt.Println("invalid department count")

		return
	}

	for range departmentCount {
		var employeesCount int
		if _, err := fmt.Scan(&employeesCount); err != nil {
			fmt.Println("invalid employees count")

			return
		}

		manager := NewTempManager()

		for range employeesCount {
			var (
				operatorSign string
				value        int
			)

			if _, err := fmt.Scan(&operatorSign, &value); err != nil {
				fmt.Println("invalid constraint")

				return
			}

			if err := manager.Apply(operatorSign, value); err != nil {
				fmt.Println("invalid operation")

				return
			}

			current, err := manager.Current()
			if err != nil {
				fmt.Println(-1)
			} else {
				fmt.Println(current)
			}
		}
	}
}
