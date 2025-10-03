package main

import "fmt"

const (
	MinTemp        = 15
	MaxTemp        = 30
	MinValidNumber = 1
	MaxValidNumber = 1000
)

func findOptimal(minTemperature, maxTemperature int, operator string, requestedTemperature int) (int, int) {
	if operator == ">=" {
		if requestedTemperature > minTemperature {
			minTemperature = requestedTemperature
		}
	}

	if operator == "<=" {
		if requestedTemperature < maxTemperature {
			maxTemperature = requestedTemperature
		}
	}

	return minTemperature, maxTemperature
}

func main() {
	var departments int

	if _, err := fmt.Scan(&departments); err != nil {
		return
	}

	if departments < MinValidNumber || departments > MaxValidNumber {
		return
	}

	for range departments {
		minTemperature := MinTemp
		maxTemperature := MaxTemp

		var employees int

		if _, err := fmt.Scan(&employees); err != nil {
			return
		}

		if employees < MinValidNumber || employees > MaxValidNumber {
			return
		}

		for range employees {
			var operator string
			var requestedTemperature int

			if _, err := fmt.Scan(&operator, &requestedTemperature); err != nil {
				return
			}

			minTemperature, maxTemperature = findOptimal(minTemperature, maxTemperature, operator, requestedTemperature)

			if minTemperature <= maxTemperature {
				fmt.Println(minTemperature)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
