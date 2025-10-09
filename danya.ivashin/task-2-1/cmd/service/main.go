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
	var temperatureValue int

	_, err := fmt.Scanln(&operation, &temperatureValue)
	if err != nil || temperatureValue < minTemperature || temperatureValue > maxTemperature {
		return 0, 0, errInvalidFormat
	}

	switch operation {
	case ">=":
		if currentHigh < temperatureValue {
			return -1, -1, nil
		}

		if currentLow < temperatureValue {
			currentLow = temperatureValue
		}
	case "<=":
		if currentLow > temperatureValue {
			return -1, -1, nil
		}

		if currentHigh > temperatureValue {
			currentHigh = temperatureValue
		}
	default:
		return 0, 0, errInvalidFormat
	}

	return currentLow, currentHigh, nil
}

func main() {
	var numberOfDepartments, numberOfEmployees int

	_, err := fmt.Scanln(&numberOfDepartments)
	if err != nil || numberOfDepartments < 1 || numberOfDepartments > 1000 {
		fmt.Println("invalid number of departments")

		return
	}

	for departmentIndex := 0; departmentIndex < numberOfDepartments; departmentIndex++ {
		_, err = fmt.Scanln(&numberOfEmployees)
		if err != nil || numberOfEmployees < 1 || numberOfEmployees > 1000 {
			fmt.Println("invalid number of employees")

			return
		}

		currentLow := minTemperature
		currentHigh := maxTemperature

		for employeeIndex := 0; employeeIndex < numberOfEmployees; employeeIndex++ {
			currentLow, currentHigh, err = adjustTemperature(currentLow, currentHigh)
			if err != nil {
				fmt.Println(err)

				return
			}

			if currentLow == -1 || currentHigh == -1 || currentLow > currentHigh {
				fmt.Println(-1)
			} else {
				fmt.Println(currentLow)
			}
		}
	}
}
