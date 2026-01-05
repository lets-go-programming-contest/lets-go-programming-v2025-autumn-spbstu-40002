package main

import (
	"errors"
	"fmt"
)

const (
	minAllowedTemp = 15
	maxAllowedTemp = 30
)

var errInvalidInput = errors.New("invalid temperature format")

func processConstraint(currentMin int, currentMax int) (int, int, error) {
	var (
		operator string
		value    int
	)

	_, err := fmt.Scanln(&operator, &value)
	if err != nil || value < minAllowedTemp || value > maxAllowedTemp {
		return 0, 0, errInvalidInput
	}

	switch operator {
	case ">=":
		if currentMax < value {
			return -1, -1, nil
		}
		if currentMin < value {
			currentMin = value
		}

	case "<=":
		if currentMin > value {
			return -1, -1, nil
		}
		if currentMax > value {
			currentMax = value
		}

	default:
		return 0, 0, errInvalidInput
	}

	return currentMin, currentMax, nil
}

func main() {
	var departmentsCount, workersCount uint16

	_, err := fmt.Scanln(&departmentsCount)
	if err != nil || departmentsCount == 0 || departmentsCount > 1000 {
		fmt.Println("invalid number of departments")
		return
	}

	for range departmentsCount {
		_, err = fmt.Scanln(&workersCount)
		if err != nil || workersCount > 1000 {
			fmt.Println("invalid number of employees")
			return
		}

		minTemp := minAllowedTemp
		maxTemp := maxAllowedTemp

		for range workersCount {
			minTemp, maxTemp, err = processConstraint(minTemp, maxTemp)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(minTemp)
		}
	}
}
