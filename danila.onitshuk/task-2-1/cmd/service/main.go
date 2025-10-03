package main

import "fmt"

func main() {
	var departmentsCount int
	if _, err := fmt.Scan(&departmentsCount); err != nil {
		return
	}

	for range departmentsCount {
		minTemperature := 15
		maxTemperature := 30
		var employeesCount int
		if _, err := fmt.Scan(&employeesCount); err != nil {
			return
		}

		for range employeesCount {
			var sign string
			if _, err := fmt.Scan(&sign); err != nil {
				return
			}
			var temperature int
			if _, err := fmt.Scan(&temperature); err != nil {
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
