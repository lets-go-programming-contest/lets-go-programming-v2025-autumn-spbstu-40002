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
		operator string
		temp     int
	)

	_, err := fmt.Scan(&operator, &temp)
	if err != nil {
		return 0, 0, errFormat
	}

	switch operator {
	case "<=":
		if temp < high {
			high = temp
		}
	case ">=":
		if temp > low {
			low = temp
		}
	default:
		return 0, 0, errFormat
	}

	return low, high, nil
}

func main() {
	var deptNum int
	if _, err := fmt.Scan(&deptNum); err != nil {
		return
	}

	for range deptNum {
		var emplNum int
		if _, err := fmt.Scan(&emplNum); err != nil {
			fmt.Println("invalid number of employees")
			return
		}

		low, high := minTemperature, maxTemperature

		for range emplNum {
			newLow, newHigh, err := adjustTemperature(low, high)
			if err != nil {
				fmt.Println(err)
				return
			}

			if newLow > newHigh {
				fmt.Println(-1)
			} else {
				fmt.Println(newLow)
				low, high = newLow, newHigh
			}
		}
	}
}