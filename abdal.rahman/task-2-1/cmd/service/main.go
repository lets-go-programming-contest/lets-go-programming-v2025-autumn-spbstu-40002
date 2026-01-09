package main

import (
	"errors"
	"fmt"
)

const (
	MinTemperature = 15
	MaxTemperature = 30
	MinOfRange     = 1
	MaxOfRange     = 1000
)

type TempData struct {
	min int
	max int
}

var (
	ErrIncorrectSign        = errors.New("incorrect sign")
	ErrTempOutOfRange       = errors.New("temperature out of range")
	ErrIncorrectDepartments = errors.New("incorrect amount of departments")
	ErrIncorrectEmployees   = errors.New("incorrect amount of employees")
	ErrDepOutOfRange        = errors.New("departments out of range")
	ErrEmpOutOfRange        = errors.New("employees out of range")
	ErrIncorrectTemp        = errors.New("incorrect temperature")
)

func newTempData(maxTemp, minTemp int) *TempData {
	return &TempData{
		min: minTemp,
		max: maxTemp,
	}
}

func (t *TempData) optimalTemp() int {
	if t.min > t.max {
		return -1
	}

	return t.min
}

func (t *TempData) adjustTemp(operator string, temp int) error {
	if temp < MinTemperature || temp > MaxTemperature {
		return ErrTempOutOfRange
	}

	switch operator {
	case ">=":
		if temp > t.min {
			t.min = temp
		}
	case "<=":
		if temp < t.max {
			t.max = temp
		}
	default:
		return ErrIncorrectSign
	}

	return nil
}

func main() {
	var departments int

	_, err := fmt.Scan(&departments)
	if err != nil {
		fmt.Println(ErrIncorrectDepartments)
		return
	}

	if departments < MinOfRange || departments > MaxOfRange {
		fmt.Println(ErrDepOutOfRange)
		return
	}

	for d := 0; d < departments; d++ {
		var employees int

		_, err = fmt.Scan(&employees)
		if err != nil {
			fmt.Println(ErrIncorrectEmployees)
			return
		}

		if employees < MinOfRange || employees > MaxOfRange {
			fmt.Println(ErrEmpOutOfRange)
			return
		}

		tempRange := newTempData(MaxTemperature, MinTemperature)

		for e := 0; e < employees; e++ {
			var operator string
			var temp int

			_, err = fmt.Scan(&operator, &temp)
			if err != nil {
				fmt.Println(ErrIncorrectTemp)
				return
			}

			err = tempRange.adjustTemp(operator, temp)
			if err != nil {
				fmt.Println(-1)
				fmt.Println()
				continue
			}

			fmt.Println(tempRange.optimalTemp())
			fmt.Println()
		}
	}
}
