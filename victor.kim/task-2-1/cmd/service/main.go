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

func adjustTemperature(currentLow int, currentHigh int) (int, int, error) {
	var (
		operator string
		value    int
	)

	_, err := fmt.Scan(&operator, &value)
	if err != nil || value < minTemperature || value > maxTemperature {
		return 0, 0, errFormat
	}

	switch operator {
	case ">=":
		if value > currentHigh {
			return -1, -1, nil
		}

		if value > currentLow {
			currentLow = value
		}
	case "<=":
		if value < currentLow {
			return -1, -1, nil
		}

		if value < currentHigh {
			currentHigh = value
		}
	default:
		return 0, 0, errFormat
	}

	return currentLow, currentHigh, nil
}

func main() {
	var departments uint16
	if _, err := fmt.Scan(&departments); err != nil || departments == 0 || departments > 1000 {
		fmt.Println("invalid number of departments")
		return
	}

	for d := uint16(0); d < departments; d++ {
		var employees uint16
		if _, err := fmt.Scan(&employees); err != nil || employees == 0 || employees > 1000 {
			fmt.Println("invalid number of employees")
			return
		}

		currentLow := minTemperature
		currentHigh := maxTemperature

		for e := uint16(0); e < employees; e++ {
			currentLow, currentHigh, err := adjustTemperature(currentLow, currentHigh)
			if err != nil {
				fmt.Println(err)
				return
			}

			if currentLow == -1 && currentHigh == -1 {
				fmt.Println(-1)
			} else {
				fmt.Println(currentLow)
			}
		}
	}
}
