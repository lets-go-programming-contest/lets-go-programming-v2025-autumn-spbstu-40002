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

	for range departmentsCount {
		minTemperature := 15
		maxTemperature := 30

		for range employeesPerDepartment {
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

			if minTemperature <= maxTemperature && minTemperature >= 15 && maxTemperature <= 30 {
				fmt.Println(minTemperature)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
