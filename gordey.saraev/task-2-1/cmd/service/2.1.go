package main

import (
	"fmt"
)

func processDepartment(employees int) error {
	minTemp := 15
	maxTemp := 30

	for range employees {
		var operator string
		var temperature int
		_, err := fmt.Scanln(&operator, &temperature)
		if err != nil {
			return err
		}

		if operator != ">=" && operator != "<=" {
			return fmt.Errorf("invalid operator: %s", operator)
		}

		if temperature < 15 || temperature > 30 {
			return fmt.Errorf("temperature %d is out of range [15, 30]", temperature)
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
	return nil
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

		err = processDepartment(employees)
		if err != nil {
			fmt.Println("Error processing department:", err)
			return
		}
	}
}
