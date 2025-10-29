package main

import (
	"errors"
	"fmt"
)

const (
	MinTemp        = 15
	MaxTemp        = 30
	MinValidNumber = 1
	MaxValidNumber = 1000
)

var (
	ErrIncorrectOperator     = errors.New("incorrect operator")
	ErrTemperatureOutOfRange = errors.New("temperature is out of range")
	ErrDepartmentsScan       = errors.New("departments could not be read")
	ErrDepartmentsOutOfRange = errors.New("departments out of range")
	ErrEmployeesScan         = errors.New("amount of employees could not be read")
	ErrEmployeesOutOfRange   = errors.New("amount of employees out of range")
	ErrTemperatureInput      = errors.New("temperature or operator could not be read")
)

type TemperatureRange struct {
	Min int
	Max int
}

func (t *TemperatureRange) Update(operator string, requested int) error {
	switch operator {
	case ">=":
		if requested > t.Min {
			t.Min = requested
		}
	case "<=":
		if requested < t.Max {
			t.Max = requested
		}
	default:
		return ErrIncorrectOperator
	}

	return nil
}

func (t *TemperatureRange) Get() (int, error) {
	if t.Min > t.Max {
		return -1, ErrTemperatureOutOfRange
	}

	return t.Min, nil
}

func main() {
	var departments int

	if _, err := fmt.Scan(&departments); err != nil {
		fmt.Println("Error:", ErrDepartmentsScan)

		return
	}

	if departments < MinValidNumber || departments > MaxValidNumber {
		fmt.Println("Error:", ErrDepartmentsOutOfRange)

		return
	}

	for range departments {
		temp := TemperatureRange{Min: MinTemp, Max: MaxTemp}

		var employees int

		if _, err := fmt.Scan(&employees); err != nil {
			fmt.Println("Error:", ErrEmployeesScan)

			return
		}

		if employees < MinValidNumber || employees > MaxValidNumber {
			fmt.Println("Error:", ErrEmployeesOutOfRange)

			return
		}

		for range employees {
			var operator string

			var requestedTemperature int

			if _, err := fmt.Scan(&operator, &requestedTemperature); err != nil {
				fmt.Println("Error:", ErrTemperatureInput)

				return
			}

			if err := temp.Update(operator, requestedTemperature); err != nil {
				return
			}

			if val, err := temp.Get(); err == nil {
				fmt.Println(val)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
