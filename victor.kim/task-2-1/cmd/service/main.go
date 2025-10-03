package main

import (
	"errors"
	"fmt"
)

const (
	minTemperature = 15
	maxTemperature = 30
	maxDepartments = 1000
	maxEmployees   = 1000
)

var ErrFormat = errors.New("invalid temperature format")

func adjustTemperature(low int, high int) (int, int, error) {
	var operation string
	var newTemp int

	_, err := fmt.Scanln(&operation, &newTemp)
	if err != nil || newTemp < minTemperature || newTemp > maxTemperature {
		return 0, 0, ErrFormat
	}

	switch operation {
	case ">=":
		if high < newTemp {
			return -1, -1, nil
		}

		if low < newTemp {
			low = newTemp
		}

	case "<=":
		if low > newTemp {
			return -1, -1, nil
		}

		if high > newTemp {
			high = newTemp
		}

	default:
		return 0, 0, ErrFormat
	}

	return low, high, nil
}

func main() {
	var departments, employees uint16

	_, err := fmt.Scanln(&departments)
	if err != nil || departments == 0 || departments > maxDepartments {
		fmt.Println("invalid number of departments")
		return
	}

	for d := 0; d < int(departments); d++ {
		_, err = fmt.Scanln(&employees)
		if err != nil || employees == 0 || employees > maxEmployees {
			fmt.Println("invalid number of employees")
			return
		}

		lowTemp := minTemperature
		highTemp := maxTemperature

		for e := 0; e < int(employees); e++ {
			lowTemp, highTemp, err = adjustTemperature(lowTemp, highTemp)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(lowTemp)
		}
	}
}
