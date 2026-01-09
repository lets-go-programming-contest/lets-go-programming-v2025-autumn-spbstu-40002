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

type TemperatureRange struct {
	Min int
	Max int
}

var (
	ErrIncorrectSign         = errors.New("incorrect sign")
	ErrTempOutOfRange        = errors.New("temperature out of range")
	ErrDepartmentsOutOfRange = errors.New("departments out of range")
	ErrEmployeesOutOfRange   = errors.New("employees out of range")
	ErrReadingInput          = errors.New("error reading input")
)

func NewTemperatureRange() *TemperatureRange {
	return &TemperatureRange{
		Min: MinTemperature,
		Max: MaxTemperature,
	}
}

func (tr *TemperatureRange) Adjust(operator string, value int) error {
	if value < MinTemperature || value > MaxTemperature {
		return ErrTempOutOfRange
	}

	switch operator {
	case ">=":
		if value > tr.Min {
			tr.Min = value
		}
	case "<=":
		if value < tr.Max {
			tr.Max = value
		}
	default:
		return ErrIncorrectSign
	}

	return nil
}

func (tr *TemperatureRange) Optimal() int {
	if tr.Min > tr.Max {
		return -1
	}
	return tr.Min
}

func main() {
	var departments int
	if _, err := fmt.Scan(&departments); err != nil {
		fmt.Println(ErrReadingInput)
		return
	}

	if departments < MinOfRange || departments > MaxOfRange {
		fmt.Println(ErrDepartmentsOutOfRange)
		return
	}

	for departmentIndex := 0; departmentIndex < departments; departmentIndex++ {
		var employees int
		if _, err := fmt.Scan(&employees); err != nil {
			fmt.Println(ErrReadingInput)
			return
		}

		if employees < MinOfRange || employees > MaxOfRange {
			fmt.Println(ErrEmployeesOutOfRange)
			return
		}

		tempRange := NewTemperatureRange()

		for employeeIndex := 0; employeeIndex < employees; employeeIndex++ {
			var operator string
			var temperature int

			if _, err := fmt.Scan(&operator, &temperature); err != nil {
				fmt.Println(ErrReadingInput)
				return
			}

			if err := tempRange.Adjust(operator, temperature); err != nil {
				fmt.Println(-1)
			} else {
				fmt.Println(tempRange.Optimal())
			}
		}
	}
}
