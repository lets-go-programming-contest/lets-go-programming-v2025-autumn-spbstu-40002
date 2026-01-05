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
	max int
	min int
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

func newTempData(maxTemp, minTemp int) (*TempData, error) {
	if minTemp > maxTemp {
		return nil, ErrTempOutOfRange
	}

	return &TempData{
		max: maxTemp,
		min: minTemp,
	}, nil
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
	var (
		departments int
		employees   int
	)

	_, err := fmt.Scan(&departments)
	if err != nil {
		fmt.Println("Error:", ErrIncorrectDepartments)

		return
	}

	if departments < MinOfRange || departments > MaxOfRange {
		fmt.Println("Error:", ErrDepOutOfRange)

		return
	}
	for range departments {

		if _, err := fmt.Scan(&employees); err != nil {
			fmt.Println("Error:", ErrIncorrectEmployees)

			return
		}

		if employees < MinOfRange || employees > MaxOfRange {
			fmt.Println("Error:", ErrEmpOutOfRange)

			return
		}

		tempRange, err := newTempData(MaxTemperature, MinTemperature)
		if err != nil {
			fmt.Println(err)
			
			return
		}
		for range employees {
			var (
				operator string
				temp     int
			)

			if _, err := fmt.Scan(&operator, &temp); err != nil {
				fmt.Println("Error:", ErrIncorrectTemp)

				return
			}

			err = tempRange.adjustTemp(operator, temp)
			if err != nil {
				fmt.Println(-1)
			} else {
				fmt.Println(tempRange.optimalTemp())
			}
		}
	}
}
