package main

import (
	"errors"
	"fmt"
)

const (
	minTemperature = 15
	maxTemperature = 30
)

var (
	errFormat      = errors.New("invalid temperature format")
	errDepartments = errors.New("invalid number of departments")
	errEmployees   = errors.New("invalid number of employees")
)

func adjustTemperature(low, high int) (int, int, error) {
	var operator string
	var temperatureValue int

	_, scanErr := fmt.Scanln(&operator, &temperatureValue)
	if scanErr != nil || temperatureValue < minTemperature || temperatureValue > maxTemperature {
		return 0, 0, errFormat
	}

	switch operator {
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

func main() {
	var numberOfDepartments int

	_, err := fmt.Scanln(&numberOfDepartments)
	if err != nil || numberOfDepartments < 1 || numberOfDepartments > 1000 {
		fmt.Println(errDepartments)
		return
	}

	for departmentIndex := 0; departmentIndex < numberOfDepartments; departmentIndex++ {
		var numberOfEmployees int
		_, err = fmt.Scanln(&numberOfEmployees)
		if err != nil || numberOfEmployees < 1 || numberOfEmployees > 1000 {
			fmt.Println(errEmployees)
			return
		}

		low := minTemperature
		high := maxTemperature

		for employeeIndex := 0; employeeIndex < numberOfEmployees; employeeIndex++ {
			var adjustErr error
			low, high, adjustErr = adjustTemperature(low, high)
			if adjustErr != nil {
				fmt.Println(adjustErr)
				return
			}

			if low == -1 || high == -1 || low > high {
				fmt.Println(-1)
			} else {
				fmt.Println(low)
			}
		}
	}
}
