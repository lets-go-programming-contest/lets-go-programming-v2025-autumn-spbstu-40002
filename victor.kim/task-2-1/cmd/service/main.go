package main

import (
	"errors"
	"fmt"
)

const (
	minTemp = 15
	maxTemp = 30
)

var errInvalidInput = errors.New("invalid temperature input")

func processPreference(currentLow, currentHigh int) (int, int, error) {
	var op string
	var temp int

	_, err := fmt.Scan(&op, &temp)
	if err != nil || temp < minTemp || temp > maxTemp {
		return 0, 0, errInvalidInput
	}

	switch op {
	case ">=":
		if currentHigh < temp {
			return -1, -1, nil
		}
		if currentLow < temp {
			currentLow = temp
		}
	case "<=":
		if currentLow > temp {
			return -1, -1, nil
		}
		if currentHigh > temp {
			currentHigh = temp
		}
	default:
		return 0, 0, errInvalidInput
	}

	return currentLow, currentHigh, nil
}

func main() {
	var numDepartments int
	_, err := fmt.Scan(&numDepartments)
	if err != nil || numDepartments <= 0 || numDepartments > 1000 {
		fmt.Println("invalid number of departments")
		return
	}

	for dep := 0; dep < numDepartments; dep++ {
		var numEmployees int
		_, err := fmt.Scan(&numEmployees)
		if err != nil || numEmployees <= 0 || numEmployees > 1000 {
			fmt.Println("invalid number of employees")
			return
		}

		low := minTemp
		high := maxTemp

		for emp := 0; emp < numEmployees; emp++ {
			low, high, err = processPreference(low, high)
			if err != nil {
				fmt.Println(err)
				return
			}

			if low == -1 && high == -1 {
				fmt.Println(-1)
				continue
			}

			fmt.Println(low)
		}
	}
}
