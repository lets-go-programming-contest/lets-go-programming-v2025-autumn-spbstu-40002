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

func adjustTemp(minTemp, maxTemp int, operator string, temp int) (int, int) {
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

	for range departments {
		minTemp := MinTempConst
		maxTemp := MaxTempConst
		valid := true

		var workers int
		_, err = fmt.Scan(&workers)
		if err != nil {
			return
		}

		if !isDataValid(workers) {
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

			if temp < MinTempConst || temp > MaxTempConst {
				valid = false
				continue
			}

			minTemp, maxTemp = adjustTemp(minTemp, maxTemp, operator, temp)

			if minTemp > maxTemp {
				valid = false
			}
		}

		if valid && minTemp <= maxTemp {
			fmt.Println(minTemp)
		} else {
			fmt.Println(-1)
		}
	}
}
