package main

import (
	"fmt"
)

func main() {
	var departmentsCount int
	fmt.Scan(&departmentsCount)

	for departmentNumber := 0; departmentNumber < departmentsCount; departmentNumber++ {
		var (
			employeesCount int
			sing           string
			temperature    int
			minTemperature int = 15
			maxTemperature int = 30
		)
		fmt.Scan(&employeesCount)

		for employee := 0; employee < employeesCount; employee++ {
			fmt.Scan(&sing, &temperature)

			switch sing {
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
