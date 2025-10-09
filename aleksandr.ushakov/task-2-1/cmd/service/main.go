package main

import (
	"errors"
	"fmt"
)

const (
	lowerTemperatureLimit = 15
	upperTemperatureLimit = 30
	minNumber             = 1
	maxNumber             = 1000
)

var errFormat = errors.New("Error")

func checkLimits(value int, minLimit int, maxLimit int) bool {
	if value >= minLimit && value <= maxLimit {
		return true
	}

	return false
}

type DepartmentTemperature struct {
	lowerLimit, upperLimit int
}

func (temperature *DepartmentTemperature) GetOptimalTemperature() int {
	if temperature.lowerLimit > temperature.upperLimit {
		return -1
	}

	return temperature.lowerLimit
}

func (temperature *DepartmentTemperature) SetOptimalTemperature(sign string, newTemp int) error {
	switch sign {
	case ">=":
		if newTemp > temperature.lowerLimit {
			temperature.lowerLimit = newTemp
		}
	case "<=":
		if newTemp < temperature.upperLimit {
			temperature.upperLimit = newTemp
		}
	default:
		return errFormat
	}

	return nil
}

func main() {
	var numberOfDepartments, numberOfPeople int
	_, err := fmt.Scanln(&numberOfDepartments)

	if err != nil || !checkLimits(numberOfDepartments, minNumber, maxNumber) {
		return
	}

	for range numberOfDepartments {
		_, err = fmt.Scanln(&numberOfPeople)
		if err != nil || !checkLimits(numberOfPeople, minNumber, maxNumber) {
			return
		}

		departmentTemp := DepartmentTemperature{lowerTemperatureLimit, upperTemperatureLimit}

		for range numberOfPeople {
			var sign string

			var temperature int

			_, err = fmt.Scan(&sign, &temperature)

			if err != nil || !checkLimits(temperature, lowerTemperatureLimit, upperTemperatureLimit) {
				return
			}

			err = departmentTemp.SetOptimalTemperature(sign, temperature)
			if err != nil {
				return
			}

			fmt.Println(departmentTemp.GetOptimalTemperature())
		}
	}
}
