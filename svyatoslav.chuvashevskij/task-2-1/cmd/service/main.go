package main

import (
	"fmt"
)

func main() {
	const maxTemperature = 30
	const minTemperature = 15

	numberOfDepartments := 0
	numberOfEmployees := 0

	_, err := fmt.Scan(&numberOfDepartments)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for i := 0; i < numberOfDepartments; i++ {
		_, err = fmt.Scan(&numberOfEmployees)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		conditions := [2]int{minTemperature, maxTemperature}
		for j := 0; j < numberOfEmployees; j++ {
			var option string
			var temperature int

			_, err = fmt.Scan(&option, &temperature)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			switch option {
			case ">=":
				if temperature > conditions[0] {
					conditions[0] = temperature
				}
			case "<=":
				if temperature < conditions[1] {
					conditions[1] = temperature
				}
			default:
				fmt.Println("Error: invalid option")
				return
			}

			if conditions[0] > conditions[1] {
				fmt.Println(-1)
			} else {
				fmt.Println(conditions[0])
			}
		}
	}
}
