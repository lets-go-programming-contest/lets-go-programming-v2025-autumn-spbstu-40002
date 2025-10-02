package main

import (
	"fmt"
	"slices"

	"lab-2-1/internal/errors"
)

func fillTemperatureTable(temperatures map[int]int, operator string, temp int) error {
	// determine acceptable temperature ranges for employees
	switch operator {
	case ">=":
		for currentTemp := temp; currentTemp <= 30; currentTemp++ {
			temperatures[currentTemp] += 1
		}
	case "<=":
		for currentTemp := temp; currentTemp >= 15; currentTemp-- {
			temperatures[currentTemp] += 1
		}
	default:
		return errors.ErrIncorrectOperator
	}
	
	return nil
}

func getAcceptableTemp(temperatures map[int]int, employeeCount int) int {
	acceptableTemperatures := make([]int, 0)

	// determine a list of temperatures suitable for each employee
	for temp := 15; temp <= 30; temp++ {
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
	var departmentCount, employeeCount, temp int
	var operator string

	// get the number of departments
	_, err := fmt.Scan(&departmentCount)
	if err != nil || !(1 <= departmentCount && departmentCount <= 1000) {
		fmt.Println(errors.ErrIncorrectDepartmentsCount)

		return
	}

	for departmentId := range departmentCount {

		// get the number of employees in the department
		_, err = fmt.Scan(&employeeCount)
		if err != nil || !(1 <= employeeCount && employeeCount <= 1000) {
			fmt.Println(errors.ErrIncorrectEmployeeCount)

			return
		}

		// initialize the map of temperatures
		temperatures := make(map[int]int)

		for employee := range employeeCount {

			// get the permissible temperature
			_, err = fmt.Scan(&operator, &temp)
			if err != nil || !(15 <= temp && temp <= 30) {
				fmt.Println(errors.ErrIncorrectTemperature)

				return
			}

			// filling out the temperature table
			err = fillTemperatureTable(temperatures, operator, temp)
			if err != nil {
				fmt.Println(err)

				return
			}

			// derive the permissible temperature
			fmt.Println(getAcceptableTemp(temperatures, employee + 1))
		}
	}
}
