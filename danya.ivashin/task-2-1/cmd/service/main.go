package main

import (
	"errors"
	"fmt"
)

const (
	minTemperature = 15
	maxTemperature = 30
)

var errFormat = errors.New("invalid temperature format")

func adjustTemperature(low int, high int) (int, int, error) {
	var operation string
	var value int

	_, err := fmt.Scanln(&operation, &value)
	if err != nil || value < minTemperature || value > maxTemperature {
		return 0, 0, errFormat
	}

	switch operation {
	case ">=":
		if high < value {
			return -1, -1, nil
		}
		if low < value {
			low = value
		}
	case "<=":
		if low > value {
			return -1, -1, nil
		}
		if high > value {
			high = value
		}
	default:
		return 0, 0, errFormat
	}

	return low, high, nil
}

func processDepartment(employees int) {
	low := minTemperature
	high := maxTemperature

	for i := 0; i < employees; i++ {
		low, high, err := adjustTemperature(low, high)
		if err != nil {
			fmt.Println(err)

			return
		}

		if low == -1 || high == -1 || low > high {
			fmt.Println(-1)
		} else {
			fmt.Println(low)
		}
	}
}

func main() {
	var departments, employees int

	_, err := fmt.Scanln(&departments)
	if err != nil || departments < 1 || departments > 1000 {
		fmt.Println("invalid number of departments")

		return
	}

	for i := 0; i < departments; i++ {
		_, err := fmt.Scanln(&employees)
		if err != nil || employees < 1 || employees > 1000 {
			fmt.Println("invalid number of employees")

			return
		}

		processDepartment(employees)
	}
}
