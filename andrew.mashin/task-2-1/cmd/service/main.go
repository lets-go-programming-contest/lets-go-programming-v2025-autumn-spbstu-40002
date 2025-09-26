package main

import (
	"fmt"
)

func main() {
	var departments int

	_, err := fmt.Scan(&departments)
	if err != nil {
		return
	}

	for range departments {
		minTemp := 15
		maxTemp := 30

		var workers int

		_, err = fmt.Scan(&workers)
		if err != nil {
			return
		}

		for range workers {
			var operator string

			_, err = fmt.Scan(&operator)
			if err != nil {
				return
			}

			var temp int

			_, err = fmt.Scan(&temp)
			if err != nil {
				return
			}

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
	}
}
