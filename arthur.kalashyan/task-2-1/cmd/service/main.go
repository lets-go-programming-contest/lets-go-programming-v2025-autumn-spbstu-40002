package main

import "fmt"

func main() {
	var departmentsCount, employeesPerDepartment int
	fmt.Scan(&departmentsCount)
	fmt.Scan(&employeesPerDepartment)

	totalEmployees := departmentsCount * employeesPerDepartment
	minTemperature, maxTemperature := 15, 30

	for i := 0; i < totalEmployees; i++ {
		var condition string
		var requestedTemperature int
		fmt.Scan(&condition, &requestedTemperature)

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
