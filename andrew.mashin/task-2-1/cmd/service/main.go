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

func findTheOptimalTemp(massiveOfTemp [2]int, operator string, temp int) [2]int {
	switch operator {
	case ">=":
		if temp > massiveOfTemp[0] {
			massiveOfTemp[0] = temp
		}
	case "<=":
		if temp < massiveOfTemp[1] {
			massiveOfTemp[1] = temp
		}
	}

	return massiveOfTemp
}

func main() {
	var departments int

	_, err := fmt.Scan(&departments)
	if err != nil || !isDataValidity(departments) {
		return
	}

	for range departments {
		massiveOfTemp := [2]int{MinTempConst, MaxTempConst}

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

			massiveOfTemp = findTheOptimalTemp(massiveOfTemp, operator, temp)

			if massiveOfTemp[0] <= massiveOfTemp[1] {
				fmt.Println(massiveOfTemp[0])
			} else {
				fmt.Println(-1)
			}
		}
	}
}
