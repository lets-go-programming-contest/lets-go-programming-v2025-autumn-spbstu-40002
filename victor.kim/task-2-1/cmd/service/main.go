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
var errContradiction = errors.New("temperature requirements cannot be satisfied")

func adjustTemperature(low int, high int) (int, int, error) {
	var (
		operation string
		newTemp   int
	)

	_, err := fmt.Scanln(&operation, &newTemp)
	if err != nil || newTemp < minTemperature || newTemp > maxTemperature {
		return low, high, errFormat
	}

	if low > high {
		return low, high, errContradiction
	}

	switch operation {
	case ">=":
		if high < newTemp {
			return low, high, errContradiction
		}

		if low < newTemp {
			low = newTemp
		}
	case "<=":
		if low > newTemp {
			return low, high, errContradiction
		}

		if high > newTemp {
			high = newTemp
		}
	default:
		return low, high, errFormat
	}

	if low > high {
		return low, high, errContradiction
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
		contradiction := false

		for range emp {
			if contradiction {
				var operation string
				var temp int
				fmt.Scanln(&operation, &temp)
				fmt.Println(-1)
				continue
			}

			var err error
			lowTemp, highTemp, err = adjustTemperature(lowTemp, highTemp)
			if err != nil {
				if err == errContradiction {
					fmt.Println(-1)
					contradiction = true
				} else {
					fmt.Println(err)

					return
				}
			} else {
				fmt.Println(lowTemp)
			}
		}
	}
}
