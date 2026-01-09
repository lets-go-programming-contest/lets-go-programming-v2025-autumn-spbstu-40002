package main

import (
	"errors"
	"fmt"
)

const (
	TempMinAllowed = 15
	TempMaxAllowed = 30
	RangeMin       = 1
	RangeMax       = 1000
)

type Temperature struct {
	MaxTemp int
	MinTemp int
}

var (
	ErrInvalidOperator     = errors.New("invalid operator")
	ErrTempOutOfBounds     = errors.New("temperature out of allowed range")
	ErrInvalidDepartments  = errors.New("invalid number of departments")
	ErrInvalidEmployees    = errors.New("invalid number of employees")
	ErrDepartmentsOutRange = errors.New("departments out of range")
	ErrEmployeesOutRange   = errors.New("employees out of range")
	ErrInvalidTemperature  = errors.New("invalid temperature input")
)

func NewTemperature(maxTemp, minTemp int) (*Temperature, error) {

	if minTemp > maxTemp {
		return nil, ErrTempOutOfBounds
	}

	return &Temperature{MaxTemp: maxTemp, MinTemp: minTemp}, nil
}

func (t *Temperature) Optimal() int {

	return t.MinTemp
}

func (t *Temperature) Adjust(operator string, tempValue int) error {

	if tempValue < TempMinAllowed || tempValue > TempMaxAllowed {
		return ErrTempOutOfBounds
	}

	switch operator {
	case ">=":
		if tempValue > t.MinTemp {
			t.MinTemp = tempValue
		}
	case "<=":
		if tempValue < t.MaxTemp {
			t.MaxTemp = tempValue
		}
	default:
		return ErrInvalidOperator
	}

	if t.MinTemp > t.MaxTemp {
		return ErrTempOutOfBounds
	}

	return nil
}

func main() {
	var numDepartments int

	if _, err := fmt.Scan(&numDepartments); err != nil {
		fmt.Println("Error:", ErrInvalidDepartments)

		return
	}

	if numDepartments < RangeMin || numDepartments > RangeMax {
		fmt.Println("Error:", ErrDepartmentsOutRange)

		return
	}

	for _ = range make([]struct{}, numDepartments) {
		var numEmployees int

		if _, err := fmt.Scan(&numEmployees); err != nil {
			fmt.Println("Error:", ErrInvalidEmployees)

			return
		}

		if numEmployees < RangeMin || numEmployees > RangeMax {
			fmt.Println("Error:", ErrEmployeesOutRange)

			return
		}

		tempData, err := NewTemperature(TempMaxAllowed, TempMinAllowed)
		if err != nil {
			fmt.Println(err)

			return
		}

		for _ = range make([]struct{}, numEmployees) {
			var operator string
			var tempValue int

			if _, err := fmt.Scan(&operator, &tempValue); err != nil {
				fmt.Println("Error:", ErrInvalidTemperature)

				return
			}

			err = tempData.Adjust(operator, tempValue)
			if err != nil {
				fmt.Println(-1)
			} else {
				fmt.Println(tempData.Optimal())
			}
		}
	}
}
