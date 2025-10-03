package main

import (
	"errors"
	"fmt"
)

const (
	minTemperature = 15
	maxTemperature = 30
)

var errInvalidFormat = errors.New("invalid temperature format")

func adjustTemperature(low int, high int) (int, int, error) {
	var op string
	var temp int

	_, scanErr := fmt.Scanln(&op, &temp)
	if scanErr != nil {
		return 0, 0, errInvalidFormat
	}

	if temp < minTemperature || temp > maxTemperature {
		return 0, 0, errInvalidFormat
	}

	if op == ">=" {
		if high < temp {
			return -1, -1, nil
		}

		if low < temp {
			low = temp
		}

	} else if op == "<=" {
		if low > temp {
			return -1, -1, nil
		}

		if high > temp {
			high = temp
		}

	} else {
		return 0, 0, errInvalidFormat
	}

	return low, high, nil
}

func main() {
	var numDepartments int
	_, err := fmt.Scan(&numDepartments)

	if err != nil {
		fmt.Println(err)
		return
	}

	if numDepartments <= 0 || numDepartments > 1000 {
		fmt.Println(-1)
		return
	}

	for range make([]struct{}, numDepartments) {
		var numEmployees int
		_, empErr := fmt.Scan(&numEmployees)

		if empErr != nil {
			fmt.Println(empErr)
			return
		}

		if numEmployees <= 0 || numEmployees > 1000 {
			fmt.Println(-1)
			return
		}

		currentLow := minTemperature
		currentHigh := maxTemperature

		for range make([]struct{}, numEmployees) {
			var adjustErr error
			currentLow, currentHigh, adjustErr = adjustTemperature(currentLow, currentHigh)

			if adjustErr != nil {
				fmt.Println(adjustErr)
				return
			}

			if currentLow == -1 || currentHigh == -1 {
				fmt.Println(-1)

				continue
			}

			fmt.Println(currentLow)
		}
	}
}
