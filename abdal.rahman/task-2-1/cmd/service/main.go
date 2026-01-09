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

// NewTemperature creates a Temperature struct ensuring Min <= Max
func NewTemperature(maxTemp, minTemp int) (*Temperature, error) {
	if minTemp > maxTemp {
		return nil, ErrTempOutOfBounds
	}

	return &Temperature{
		MaxTemp: maxTemp,
		MinTemp: minTemp,
	}, nil
}

// Optimal returns the current optimal temperature (MinTemp in this task)
func (t *Temperature) Optimal() int {
	return t.MinTemp
}

// Adjust updates MinTemp or MaxTemp based on operator and temperature input
func (t *Temperature) Adjust(op string, temp int) error {
	if temp < TempMinAllowed || temp > TempMaxAllowed {
		return ErrTempOutOfBounds
	}

	switch op {
	case ">=":
		if temp > t.MinTemp {
			t.MinTemp = temp
		}
	case "<=":
		if temp < t.MaxTemp {
			t.MaxTemp = temp
		}
	default:
		return ErrInvalidOperator
	}

	// Ensure MinTemp does not exceed MaxTemp
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

	for d := 0; d < numDepartments; d++ {
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

		for e := 0; e < numEmployees; e++ {
			var op string
			var temp int
			if _, err := fmt.Scan(&op, &temp); err != nil {
				fmt.Println("Error:", ErrInvalidTemperature)
				return
			}

			err = tempData.Adjust(op, temp)
			if err != nil {
				fmt.Println(-1)
			} else {
				fmt.Println(tempData.Optimal())
			}
		}
	}
}
