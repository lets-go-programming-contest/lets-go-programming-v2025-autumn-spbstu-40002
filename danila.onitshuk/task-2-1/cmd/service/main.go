package main

import (
	"fmt"
)

func main() {
	var departmentsCount int
	if _, err := fmt.Scan(&departmentsCount); err != nil {
		return
	}

	for departmentNumber := 0; departmentNumber < departmentsCount; departmentNumber++ {
		var (
			employeesCount int
			sign           string
			temperature    int
			minTemperature int = 15
			maxTemperature int = 30
		)
		if _, err := fmt.Scan(&employeesCount); err != nil {
			return
		}

		for employee := 0; employee < employeesCount; employee++ {
			if _, err := fmt.Scan(&sign, &temperature); err != nil {
				return
			}

			switch sign {
			case ">=":
				if temperature > minTemperature {
					minTemperature = temperature
				}
			case "<=":
				if temperature < maxTemperature {
					maxTemperature = temperature
				}
			}

			if minTemperature <= maxTemperature && minTemperature >= 15 && maxTemperature <= 30 {
				fmt.Println(minTemperature)
			} else {
				fmt.Println(-1)
			}

		}
	}
}
