package main

import "fmt"

func main() {
	var departmentsCount, employeesPerDepartment int
	if _, err := fmt.Scan(&departmentsCount); err != nil {
		return
	}
	if _, err := fmt.Scan(&employeesPerDepartment); err != nil {
		return
	}

	totalEmployees := departmentsCount * employeesPerDepartment
	minTemperature, maxTemperature := 15, 30

	for i := 0; i < totalEmployees; i++ {
		var condition string
		var requestedTemperature int

		if _, err := fmt.Scan(&condition, &requestedTemperature); err != nil {
			return
		}

		if condition == ">=" && requestedTemperature > minTemperature {
			minTemperature = requestedTemperature
		}

		if condition == "<=" && requestedTemperature < maxTemperature {
			maxTemperature = requestedTemperature
		}

		if minTemperature <= maxTemperature {
			fmt.Println(minTemperature)
		} else {
			fmt.Println(-1)
		}
	}
}
