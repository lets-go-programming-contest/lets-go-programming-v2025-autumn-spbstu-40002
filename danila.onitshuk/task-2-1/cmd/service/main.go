package main

import "fmt"

func department() {
	minTemperature := 15
	maxTemperature := 30

	var employeesCount int

	_, err := fmt.Scan(&employeesCount)
	if err != nil {
		return
	}

	for range employeesCount {
		var (
			sign        string
			temperature int
		)

		_, err := fmt.Scan(&sign, &temperature)
		if err != nil {
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

func main() {
	var departmentsCount int

	_, err := fmt.Scan(&departmentsCount)
	if err != nil {
		return
	}

	for range departmentsCount {
		department()
	}
}
