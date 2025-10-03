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
	var operator string
	var value int

	_, err := fmt.Scan(&operator, &value)
	if err != nil {
		return 0, 0, errFormat
	}

	if value < minTemperature || value > maxTemperature {
		return 0, 0, errFormat
	}

	switch operator {
	case ">=":
		if value > high {
			return -1, -1, nil
		}

		if value > low {
			low = value
		}
	case "<=":
		if value < low {
			return -1, -1, nil
		}

		if value < high {
			high = value
		}
	default:
		return 0, 0, errFormat
	}

	return low, high, nil
}

func main() {
	var departments int
	n, err := fmt.Scan(&departments)
	if n != 1 || err != nil {
		return
	}

	if departments <= 0 || departments > 1000 {
		fmt.Println("invalid number of departments")
		return
	}

	for d := 0; d < departments; d++ {
		var employees int
		n, err = fmt.Scan(&employees)
		if n != 1 || err != nil {
			return
		}

		if employees <= 0 || employees > 1000 {
			fmt.Println("invalid number of employees")
			return
		}

		lowTemp := minTemperature
		highTemp := maxTemperature

		for e := 0; e < employees; e++ {
			var err error
			lowTemp, highTemp, err = adjustTemperature(lowTemp, highTemp)
			if err != nil {
				fmt.Println(err)
				return
			}

			if lowTemp == -1 && highTemp == -1 {
				fmt.Println(-1)
			} else {
				fmt.Println(lowTemp)
			}
		}
	}
}
