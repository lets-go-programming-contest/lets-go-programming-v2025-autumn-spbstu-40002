package main

import "fmt"

const (
	MinTempConst   = 15
	MaxTempConst   = 30
	MinCorrectData = 1
	MaxCorrectData = 1000
)

func isDataValid(value int) bool {
	return value >= MinCorrectData && value <= MaxCorrectData
}

func adjustTemperature(minTemp, maxTemp int, operator string, temperature int) (int, int) {
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

	for range make([]struct{}, departments) {
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

		for range make([]struct{}, workers) {
			var operator string
			_, err = fmt.Scan(&operator)
			if err != nil {
				return
			}

			var temperature int
			_, err = fmt.Scan(&temperature)
			if err != nil {
				return
			}

			if temperature < MinTempConst || temperature > MaxTempConst {
				fmt.Println(-1)
				continue
			}

			minTemp, maxTemp = adjustTemperature(minTemp, maxTemp, operator, temperature)

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
