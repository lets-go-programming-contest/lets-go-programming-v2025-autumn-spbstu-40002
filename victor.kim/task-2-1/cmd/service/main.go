package main

import (
	"errors"
	"fmt"
)

const (
	MinTemperature = 15
	MaxTemperature = 30
)

var ErrFormat = errors.New("invalid temperature format")
var ErrContradiction = errors.New("temperature requirements cannot be satisfied")

func adjustTemperature(low int, high int) (int, int, error) {
	var operation string
	var newTemp int

	_, err := fmt.Scanln(&operation, &newTemp)
	if err != nil || newTemp < MinTemperature || newTemp > MaxTemperature {
		return low, high, ErrFormat
	}

	if low > high {
		return low, high, ErrContradiction
	}

	switch operation {
	case ">=":
		if high < newTemp {
			return low, high, ErrContradiction
		}

		if low < newTemp {
			low = newTemp
		}
	case "<=":
		if low > newTemp {
			return low, high, ErrContradiction
		}

		if high > newTemp {
			high = newTemp
		}
	default:
		return low, high, ErrFormat
	}

	if low > high {
		return low, high, ErrContradiction
	}

	return low, high, nil
}

func readUint16(prompt string, maxValue uint16) (uint16, error) {
	var value uint16
	_, err := fmt.Scanln(&value)
	if err != nil || value == 0 || value > maxValue {
		return 0, errors.New(prompt)
	}

	return value, nil
}

func main() {
	dep, err := readUint16("invalid number of departments", 1000)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < int(dep); i++ {
		emp, err := readUint16("invalid number of employees", 1000)
		if err != nil {
			fmt.Println(err)
			return
		}

		lowTemp := MinTemperature
		highTemp := MaxTemperature
		contradiction := false

		for j := 0; j < int(emp); j++ {
			if contradiction {
				var op string
				var temp int
				_, _ = fmt.Scanln(&op, &temp)
				fmt.Println(-1)
				continue
			}

			lowTemp, highTemp, err = adjustTemperature(lowTemp, highTemp)
			if err != nil {
				if errors.Is(err, ErrContradiction) {
					fmt.Println(-1)
					contradiction = true
				} else {
					fmt.Println(err)
					return
				}
			} else {
				fmt.Println(lowTemp)
			}
		}
	}
}
