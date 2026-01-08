package main

import (
	"errors"
	"fmt"
)

const (
	MinAllowedTemp = 15
	MaxAllowedTemp = 30
	MinDepartments = 1
	MaxDepartments = 1000
	MinEmployees   = 1
	MaxEmployees   = 1000
)

type TemperatureRange struct {
	maximum int
	minimum int
}

var (
	ErrInvalidOperator     = errors.New("incorrect sign")
	ErrTempOutOfBounds     = errors.New("temperature out of range")
	ErrInvalidDeptCount    = errors.New("incorrect amount of departments")
	ErrInvalidEmpCount     = errors.New("incorrect amount of employees")
	ErrDeptCountOutOfRange = errors.New("departments out of range")
	ErrEmpCountOutOfRange  = errors.New("employees out of range")
	ErrInvalidTempInput    = errors.New("incorrect temperature")
	ErrInvalidTempRange    = errors.New("incorrect temperature range")
)

func NewTemperatureRange(maxTemp, minTemp int) (*TemperatureRange, error) {
	if minTemp > maxTemp {
		return nil, ErrInvalidTempRange
	}

	return &TemperatureRange{
		maximum: maxTemp,
		minimum: minTemp,
	}, nil
}

func (tr *TemperatureRange) GetOptimalTemp() int {
	if tr.minimum > tr.maximum {
		return -1
	}

	return tr.minimum
}

func (tr *TemperatureRange) UpdateRange(operator string, temp int) error {
	if temp < MinAllowedTemp || temp > MaxAllowedTemp {
		return ErrTempOutOfBounds
	}

	switch operator {
	case ">=":
		if temp > tr.minimum {
			tr.minimum = temp
		}
	case "<=":
		if temp < tr.maximum {
			tr.maximum = temp
		}
	default:
		return ErrInvalidOperator
	}

	if tr.minimum > tr.maximum {
		return ErrInvalidTempRange
	}

	return nil
}

func ProcessDepartment() error {
	var employeeCount int

	if _, err := fmt.Scan(&employeeCount); err != nil {
		return ErrInvalidEmpCount
	}

	if employeeCount < MinEmployees || employeeCount > MaxEmployees {
		return ErrEmpCountOutOfRange
	}

	tempRange, err := NewTemperatureRange(MaxAllowedTemp, MinAllowedTemp)
	if err != nil {
		return err
	}

	for i := 0; i < employeeCount; i++ {

		var operator string
		var temperature int

		if _, err := fmt.Scan(&operator, &temperature); err != nil {
			fmt.Println(-1)
			continue
		}

		err := tempRange.UpdateRange(operator, temperature)
		if err != nil {
			fmt.Println(-1)
		} else {
			fmt.Println(tempRange.GetOptimalTemp())
		}
	}

	return nil
}

func main() {
	var departmentCount int

	_, err := fmt.Scan(&departmentCount)
	if err != nil {
		fmt.Println("Error:", ErrInvalidDeptCount)
		return
	}

	if departmentCount < MinDepartments || departmentCount > MaxDepartments {
		fmt.Println("Error:", ErrDeptCountOutOfRange)
		return
	}

	for i := 0; i < departmentCount; i++ {

		err := ProcessDepartment()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
}
