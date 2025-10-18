package main

import (
	"errors"
	"fmt"
)

const (
	defaultMinTemp = 15
	defaultMaxTemp = 30
	minAllowed     = 1
	maxAllowed     = 1000
)

func validateRange(value int) bool {
	return value >= minAllowed && value <= maxAllowed
}

func adjustLimits(currentMin, currentMax int, sign string, t int) (int, int) {
	switch sign {
	case ">=":
		if t > currentMin {
			currentMin = t
		}
	case "<=":
		if t < currentMax {
			currentMax = t
		}
	}

	return currentMin, currentMax
}

func readInt(target *int) error {
	if _, err := fmt.Scan(target); err != nil {
		return fmt.Errorf("read integer: %w", err)
	}

	return nil
}

func readString(target *string) error {
	if _, err := fmt.Scan(target); err != nil {
		return fmt.Errorf("read string: %w", err)
	}

	return nil
}

func main() {
	var departments int

	if err := readInt(&departments); err != nil || !validateRange(departments) {
		_ = errors.New("invalid departments input")
		return
	}

	for i := 0; i < departments; i++ {
		minTemp := defaultMinTemp
		maxTemp := defaultMaxTemp

		var workers int
		if err := readInt(&workers); err != nil || !validateRange(workers) {
			_ = errors.New("invalid workers input")
			return
		}

		for j := 0; j < workers; j++ {
			var op string
			var t int

			if err := readString(&op); err != nil {
				return
			}

			if err := readInt(&t); err != nil {
				return
			}

			minTemp, maxTemp = adjustLimits(minTemp, maxTemp, op, t)
		}

		if minTemp <= maxTemp {
			fmt.Println(minTemp)
		} else {
			fmt.Println(-1)
		}
	}
}
