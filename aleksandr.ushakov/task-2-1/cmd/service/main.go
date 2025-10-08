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

func checkLimits(value int, minLimit int, maxLimit int) bool {
	if value >= minLimit && value <= maxLimit {
		return true
	}
	return false
}

func getOptimalTemperature(sign string, temperature int, lowerLimit *int, upperLimit *int) (int, error) {
	switch sign {
	case ">=":
		if temperature > *lowerLimit {
			*lowerLimit = temperature
		}
	case "<=":
		if temperature < *upperLimit {
			*upperLimit = temperature
		}
	default:
		return 0, errors.New("error")
	}

	if *lowerLimit > *upperLimit {
		return -1, nil
	}
	return *lowerLimit, nil
}

func main() {
	var numberOfDepartments int
	var numberOfPeople int
	_, err := fmt.Scanln(&numberOfDepartments)
	if err != nil || !checkLimits(numberOfDepartments, minNumber, maxNumber) {
		return
	}

	for i := 0; i < numberOfDepartments; i++ {
		_, err = fmt.Scanln(&numberOfPeople)
		if err != nil || !checkLimits(numberOfPeople, minNumber, maxNumber) {
			return
		}
		lowerDepartmentLimit := 15
		upperDepartmentLimit := 30
		for j := 0; j < numberOfPeople; j++ {
			var sign string
			var temperature int
			_, err = fmt.Scan(&sign, &temperature)
			if err != nil || !checkLimits(temperature, lowerTemperatureLimit, upperTemperatureLimit) {
				return
			}
			optimalTemp, err := getOptimalTemperature(sign, temperature, &lowerDepartmentLimit, &upperDepartmentLimit)
			if err != nil {
				return
			}
			fmt.Println(optimalTemp)
		}
	}
}
