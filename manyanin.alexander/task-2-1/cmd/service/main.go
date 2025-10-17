package main

import (
	"fmt"
)

const (
	MinTemp = 15
	MaxTemp = 30
	MinData = 1
	MaxData = 1000
)

func isValid(data int) bool {
	return data >= MinData && data <= MaxData
}

func findBestTemp(minTemp int, maxTemp int, operator string, temp int) (int, int) {
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

	return minTemp, maxTemp
}

func main() {
	var otdels int

	_, err := fmt.Scan(&otdels)
	if err != nil || !isValid(otdels) {
		return
	}

	for range otdels {
		minTemp := MinTemp
		maxTemp := MaxTemp

		var workers int

		_, err := fmt.Scan(&workers)
		if err != nil || !isValid(workers) {
			return
		}

		for range workers {
			var operator string
			var temp int

			_, err = fmt.Scan(&operator)
			if err != nil {
				return
			}

			_, err = fmt.Scan(&temp)
			if err != nil {
				return
			}

			minTemp, maxTemp = findBestTemp(minTemp, maxTemp, operator, temp)

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
