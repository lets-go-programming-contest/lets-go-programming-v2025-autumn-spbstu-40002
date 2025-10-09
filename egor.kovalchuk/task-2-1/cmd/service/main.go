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
	var (
		operation string
		newTemp   int
	)

	_, err := fmt.Scanln(&operation, &newTemp)
	if err != nil || newTemp < minTemperature || newTemp > maxTemperature {
		return 0, 0, errFormat
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
		return 0, 0, errFormat
	}

	return low, high, nil
}

func main() {
	var dep, emp uint16

	_, err := fmt.Scanln(&dep)
	if err != nil || dep == 0 || dep > 1000 {
		fmt.Println("invalid number of departments")

		return
	}

	for range dep {
		_, err = fmt.Scanln(&emp)
		if err != nil || emp > 1000 {
			fmt.Println("invalid number of employees")

			return
		}

		lowTemp := minTemperature
		highTemp := maxTemperature

		for range emp {
			lowTemp, highTemp, err = adjustTemperature(lowTemp, highTemp)
			if err != nil {
				fmt.Println(err)

				return
			}

			fmt.Println(lowTemp)
		}
	}
}
