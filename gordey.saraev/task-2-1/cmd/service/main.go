package main

import (
	"fmt"
)

func processDepartment(employees int) {
	minTemp := 15
	maxTemp := 30

	for range employees {
		var operator string
		var temperature int
		_, err := fmt.Scanln(&operator, &temperature)
		if err != nil {
			fmt.Println("Error reading operator and temperature")
			return
		}

		if operator != ">=" && operator != "<=" {
			fmt.Println("Invalid operator. Must be '>=' or '<='")
			return
		}

		if temperature < 15 || temperature > 30 {
			fmt.Println("Temperature is out of range [15, 30]")
			return
		}

		switch operator {
		case ">=":
			if temperature > minTemp {
				minTemp = temperature
			}
		case "<=":
			if temperature < maxTemp {
				maxTemp = temperature
			}
		}

		if minTemp <= maxTemp {
			fmt.Println(minTemp)
		} else {
			fmt.Println(-1)
		}
	}
}

func main() {
	var departments int
	_, err := fmt.Scanln(&departments)
	if err != nil {
		fmt.Println("Error reading departments count")
		return
	}

	if departments < 1 || departments > 1000 {
		fmt.Println("Departments is out of range [1, 1000]")
		return
	}

	for range departments {
		var employees int
		_, err = fmt.Scanln(&employees)
		if err != nil {
			fmt.Println("Error reading employees count")
			return
		}

		if employees < 1 || employees > 1000 {
			fmt.Println("Employees is out of range [1, 1000]")
			return
		}

		processDepartment(employees)
	}
}
