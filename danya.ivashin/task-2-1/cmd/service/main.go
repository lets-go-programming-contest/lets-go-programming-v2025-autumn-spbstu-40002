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

func adjustTemperature(low, high int) (int, int, error) {
	var operation string
	var temperatureValue int

	_, err := fmt.Scanln(&operation, &temperatureValue)
	if err != nil || temperatureValue < minTemperature || temperatureValue > maxTemperature {
		return 0, 0, errFormat
	}

	switch operation {
	case ">=":
		if high < temperatureValue {
			return -1, -1, nil
		}
		if low < temperatureValue {
			low = temperatureValue
		}
	case "<=":
		if low > temperatureValue {
			return -1, -1, nil
		}
		if high > temperatureValue {
			high = temperatureValue
		}
	default:
		return 0, 0, errFormat
	}

	return low, high, nil
}

func processDepartment(numberOfEmployees int) error {
	if numberOfEmployees < 1 || numberOfEmployees > 1000 {
		return errors.New("invalid number of employees")
	}

	low := minTemperature
	high := maxTemperature

	for range numberOfEmployees {
		var err error
		low, high, err = adjustTemperature(low, high)
		if err != nil {
			return err
		}

		if low == -1 || high == -1 || low > high {
			fmt.Println(-1)
		} else {
			fmt.Println(low)
		}
	}

	return nil
}

func main() {
	var numberOfDepartments int

	_, err := fmt.Scanln(&numberOfDepartments)
	if err != nil || numberOfDepartments < 1 || numberOfDepartments > 1000 {
		fmt.Println("invalid number of departments")
		return
	}

	for range numberOfDepartments {
		var numberOfEmployees int

		_, err = fmt.Scanln(&numberOfEmployees)
		if err != nil {
			fmt.Println("invalid number of employees")
			return
		}

		err = processDepartment(numberOfEmployees)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
