package main

import "fmt"

const (
	MinTempConst   = 15
	MaxTempConst   = 30
	MinCorrectData = 1
	MaxCorrectData = 1000
)

func isDataValid(x int) bool {
	return x >= MinCorrectData && x <= MaxCorrectData
}

func adjustTemp(minTemp, maxTemp int, op string, temp int) (int, int) {
	switch op {
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
	if err != nil || !isDataValid(departments) {
		return
	}

	for range make([]struct{}, departments) {
		minTemp := MinTempConst
		maxTemp := MaxTempConst

		var workers int
		_, err := fmt.Scan(&workers)
		if err != nil || !isDataValid(workers) {
			return
		}

		for range make([]struct{}, workers) {
			var op string
			_, err = fmt.Scan(&op)
			if err != nil {
				return
			}

			var temp int
			_, err = fmt.Scan(&temp)
			if err != nil {
				return
			}

			if temp < MinTempConst || temp > MaxTempConst {
				fmt.Println(-1)
				continue
			}

			minTemp, maxTemp = adjustTemp(minTemp, maxTemp, op, temp)

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
