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

func adjustTemperature(currentLow int, currentHigh int) (int, int, error) {
	var operation string
	var temperature int

	_, scanErr := fmt.Scanln(&operation, &temperature)
	if scanErr != nil {
		return 0, 0, errInvalidFormat
	}

	if temperature < minTemperature || temperature > maxTemperature {
		return 0, 0, errInvalidFormat
	}

	switch operation {
	case ">=":
		if currentHigh < temperature {
			return -1, -1, nil
		}
		if currentLow < temperature {
			currentLow = temperature
		}
	case "<=":
		if currentLow > temperature {
			return -1, -1, nil
		}
		if currentHigh > temperature {
			currentHigh = temperature
		}
	default:
		return 0, 0, errInvalidFormat
	}

	return currentLow, currentHigh, nil
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

	for dep := 0; dep < numDepartments; dep++ {
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

		for emp := 0; emp < numEmployees; emp++ {
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
