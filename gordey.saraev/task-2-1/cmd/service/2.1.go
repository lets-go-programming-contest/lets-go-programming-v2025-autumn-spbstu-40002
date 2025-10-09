package main

import (
	"fmt"
)

func main() {
	var departments, employee int
	var operator string
	var temperature int

	_, err := fmt.Scan(&departments)
	if err != nil {
		fmt.Println("Invalid number of departments")
		return
	}
	if departments < 1 || departments > 1000 {
		fmt.Println("N is out of range [1, 1000]")
		return
	}

	for i := 0; i < departments; i++ {
		_, err = fmt.Scan(&employee)
		if err != nil {
			fmt.Println("Invalid number of employees")
			return
		}
		if employee < 1 || employee > 1000 {
			fmt.Println("K is out of range [1, 1000]")
			return
		}

		minTemp := 15
		maxTemp := 30

		for j := 0; j < employee; j++ {
			_, err = fmt.Scan(&operator, &temperature)
			if err != nil || temperature > 30 || temperature < 15 {
				fmt.Println("Invalid temperature constraint format")
				return
			}

			switch operator {
			case ">=":
				if temperature >= minTemp {
					minTemp = temperature
				}
			case "<=":
				if temperature <= maxTemp {
					maxTemp = temperature
				}
			default:
				fmt.Println(-1)
				continue
			}

			if minTemp > maxTemp {
				fmt.Println(-1)
			} else {
				fmt.Println(maxTemp)
			}
		}
	}
}
