package main

import "fmt"

func main() {
	var departmentCount int
	_, _ = fmt.Scan(&departmentCount)

	for i := 0; i < departmentCount; i++ {
		var employeesCount int
		_, _ = fmt.Scan(&employeesCount)

		minTemp := 15
		maxTemp := 30

		for j := 0; j < employeesCount; j++ {
			var operator string
			var temp int
			_, _ = fmt.Scan(&operator, &temp)

			switch operator {
			case ">=":
				if temp > minTemp {
					minTemp = temp
				}
			case "<=":
				if temp < maxTemp {
					maxTemp = temp
				}
			}

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
		fmt.Println()
	}
}
