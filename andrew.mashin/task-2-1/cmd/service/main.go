package main

import (
	"fmt"
)

const (
	MinTempConst = 15
	MaxTempConst = 30
)

func isDataValidity(data int) bool {
	return data >= 1 && data <= 1000
}

func main() {
	var departments int

	_, err := fmt.Scan(&departments)
	if err != nil || !isDataValidity(departments) {
		return
	}

	for range departments {
		minTemp := MinTempConst
		maxTemp := MaxTempConst

		var workers int

		_, err = fmt.Scan(&workers)
		if err != nil || !isDataValidity(workers) {
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
