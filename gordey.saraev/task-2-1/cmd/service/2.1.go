package main

import (
	"fmt"
)

func main() {
	var departments, employee int
	var operator string
	var temperature int

	_, err := fmt.Scanf("%d\n", &departments)
	if err != nil {
		fmt.Fprintln(nil, "Invalid number of departments")
		return
	}
	if departments < 1 || departments > 1000 {
		fmt.Fprintln(nil, "N is out of range [1, 1000]")
		return
	}

	for i := 0; i < departments; i++ {
		_, err = fmt.Scanf("%d\n", &employee)
		if err != nil {
			fmt.Fprintln(nil, "Invalid number of employees")
			return
		}
		if employee < 1 || employee > 1000 {
			fmt.Fprintln(nil, "K is out of range [1, 1000]")
			return
		}

		minTemp := 15
		maxTemp := 30

		for j := 0; j < employee; j++ {
			_, err = fmt.Scanf("%s %d\n", &operator, &temperature)
			if err != nil {
				fmt.Fprintln(nil, "Invalid temperature constraint format")
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
			default:
				fmt.Println(-1)
				continue
			}

			if minTemp <= maxTemp && minTemp >= 15 && maxTemp <= 30 {
				fmt.Println(maxTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
