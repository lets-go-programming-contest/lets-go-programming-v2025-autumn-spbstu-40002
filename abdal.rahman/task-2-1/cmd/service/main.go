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

func NewTempData() *TempData {
	return &TempData{
		min: MinTemperature,
		max: MaxTemperature,
	}
}

func (t *TempData) Adjust(operator string, temp int) error {
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

func (t *TempData) Result() int {
	if t.min <= t.max {
		return t.min
	}
	return -1
}

func main() {
	var departments int
	if _, err := fmt.Scan(&departments); err != nil || departments < MinOfRange || departments > MaxOfRange {
		fmt.Println(ErrIncorrectDepartments)
		return
	}

	for i := 0; i < departments; i++ {
		var employees int
		if _, err := fmt.Scan(&employees); err != nil || employees < MinOfRange || employees > MaxOfRange {
			fmt.Println(ErrIncorrectEmployees)
			return
		}

		tempRange := NewTempData()

		for j := 0; j < employees; j++ {
			var operator string
			var temp int
			if _, err := fmt.Scan(&operator, &temp); err != nil {
				fmt.Println(ErrIncorrectTemp)
				return
			}

			if err := tempRange.Adjust(operator, temp); err != nil {
				fmt.Println(-1)
			} else {
				fmt.Println(tempRange.Result())
			}
		}
	}
}
