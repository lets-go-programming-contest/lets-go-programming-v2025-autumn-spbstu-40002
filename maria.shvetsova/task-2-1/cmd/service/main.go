package main

import (
	"fmt"
)

func findOptimalTemp(sign string, temp int, maxTemp, minTemp, optimalTemp *int) {
	switch sign {
	case ">=":
		if temp > *maxTemp {
			*optimalTemp = -1
		} else if temp > *minTemp {
			*minTemp = temp
		}
	case "<=":
		if temp < *minTemp {
			*optimalTemp = -1
		} else if temp < *maxTemp {
			*maxTemp = temp
		}
	default:
		return
	}

	if *optimalTemp != -1 {
		*optimalTemp = *minTemp
	}

	fmt.Println(*optimalTemp)
}

func main() {
	var numOfDepartments int

	_, err := fmt.Scan(&numOfDepartments)
	if err != nil {
		return
	}

	if numOfDepartments < 1 || numOfDepartments > 1000 {
		return
	}

	for range numOfDepartments {
		var numOfEmployees int

		_, err = fmt.Scan(&numOfEmployees)
		if err != nil {
			return
		}

		if numOfEmployees < 1 || numOfEmployees > 1000 {
			return
		}

		var sign string

		var temperature int

		maxTemp := 30
		minTemp := 15
		optimalTemp := 15

		for range numOfEmployees {
			_, err = fmt.Scan(&sign, &temperature)
			if err != nil {
				return
			}

			findOptimalTemp(sign, temperature, &maxTemp, &minTemp, &optimalTemp)
		}
	}
}
