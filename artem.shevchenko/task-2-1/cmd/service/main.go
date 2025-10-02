package main

import (
	"errors"
	"fmt"
	"slices"
)

const (
	minTemp              = 15
	maxTemp              = 30
	minInitialConditions = 1
	maxInitialConditions = 1000
)

var (
	ErrIncorrectTemperature      = errors.New("the temperature is not between 15 and 30")
	ErrIncorrectOperator         = errors.New("the range sign can only be <= or >=")
	ErrIncorrectDepartmentsCount = errors.New("the number of departments must be between 1 and 1000")
	ErrIncorrectEmployeeCount    = errors.New("the number of employees must be between 1 and 1000")
)

func readOperatorAndTemp() (string, int, bool) {
	var (
		operator string
		temp     int
	)

	_, err := fmt.Scan(&operator, &temp)
	if err != nil || !(minTemp <= temp && temp <= maxTemp) {
		fmt.Println(ErrIncorrectTemperature)

		return "", 0, false
	}

	if operator != ">=" && operator != "<=" {
		fmt.Println(ErrIncorrectOperator)

		return "", 0, false
	}

	return operator, temp, true
}

func fillTemperatureTable(temperatures map[int]int, operator string, temp int) {
	// determine acceptable temperature ranges for employees
	switch operator {
	case ">=":
		for currentTemp := temp; currentTemp <= maxTemp; currentTemp++ {
			temperatures[currentTemp] += 1
		}
	case "<=":
		for currentTemp := temp; currentTemp >= minTemp; currentTemp-- {
			temperatures[currentTemp] += 1
		}
	}
}

func getAcceptableTemp(temperatures map[int]int, employeeCount int) int {
	acceptableTemperatures := make([]int, 0)

	// determine a list of temperatures suitable for each employee
	for temp := minTemp; temp <= maxTemp; temp++ {
		if temperatures[temp] == employeeCount {
			acceptableTemperatures = append(acceptableTemperatures, temp)
		}
	}

	// find the minimum temperature
	if len(acceptableTemperatures) != 0 {
		return slices.Min(acceptableTemperatures)
	}

	return -1
}

func main() {
	var (
		departmentCount int
		employeeCount   int
	)

	// get the number of departments
	_, err := fmt.Scan(&departmentCount)
	if err != nil || !(minInitialConditions <= departmentCount && departmentCount <= maxInitialConditions) {
		fmt.Println(ErrIncorrectDepartmentsCount)

		return
	}

	for range departmentCount {
		// get the number of employees in the department
		_, err = fmt.Scan(&employeeCount)
		if err != nil || !(minInitialConditions <= employeeCount && employeeCount <= maxInitialConditions) {
			fmt.Println(ErrIncorrectEmployeeCount)

			return
		}

		// initialize the map of temperatures
		temperatures := make(map[int]int)

		for employee := range employeeCount {
			// get the permissible temperature and operator
			operator, temp, ok := readOperatorAndTemp()
			if !ok {
				return
			}

			// filling out the temperature table
			fillTemperatureTable(temperatures, operator, temp)

			// derive the permissible temperature
			fmt.Println(getAcceptableTemp(temperatures, employee+1))
		}
	}
}
