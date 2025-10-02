package main

import (
	"fmt"
)

func findOptimalTemp(sign string, temp int, maxTemp, minTemp, optimalTemp *int) {
	switch sign {
	case ">=":
		switch {
		case temp <= *maxTemp && temp <= *minTemp || temp <= *optimalTemp:
			*minTemp = temp
		case temp <= *maxTemp && !(temp <= *minTemp || temp <= *optimalTemp):
			*optimalTemp = temp
		default:
			*optimalTemp = -1
		}
	case "<=":
		switch {
		case temp >= *minTemp && temp >= *maxTemp || temp >= *optimalTemp:
			*maxTemp = temp
		case temp <= *maxTemp && !(temp >= *maxTemp || temp >= *optimalTemp):
			*optimalTemp = temp
		default:
			*optimalTemp = temp
		}
	default:
		return
	}
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
			_, err = fmt.Scan(&sign)
			if err != nil {
				return
			}

			_, err = fmt.Scan(&temperature)
			if err != nil {
				return
			}

			findOptimalTemp(sign, temperature, &maxTemp, &minTemp, &optimalTemp)

			fmt.Println(optimalTemp)
		}
	}
}
