package main

import (
	"fmt"
)

const (
	MinTempConst   = 15
	MaxTempConst   = 30
	MinCorrectData = 1
	MaxCorrectData = 1000
)

func isDataValid(data int) bool {
	return data >= MinCorrectData && data <= MaxCorrectData
}

func findOptimalTemp(minTemp int, maxTemp int, operator string, temp int) (int, int) {
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
	var departments int

	_, err := fmt.Scan(&departments)
	if err != nil {
		return
	}

	if !isDataValid(departments) {
		return
	}

	for i := 0; i < departments; i++ {
		minTemp := MinTempConst
		maxTemp := MaxTempConst

		var workers int

		_, err = fmt.Scan(&workers)
		if err != nil {
			return
		}

		if !isDataValid(workers) {
			return
		}

		for j := 0; j < workers; j++ {
			var operator string
			_, err = fmt.Scan(&operator)
			if err != nil {
				continue
			}

			var temp int
			_, err = fmt.Scan(&temp)
			if err != nil {
				continue
			}

			if temp < MinTempConst || temp > MaxTempConst {
				fmt.Println(-1)
				continue
			}

			minTemp, maxTemp = findOptimalTemp(minTemp, maxTemp, operator, temp)

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
