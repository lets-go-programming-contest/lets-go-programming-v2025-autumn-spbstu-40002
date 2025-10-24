package main

import (
	"errors"
	"fmt"
)

const (
	minTemperatureConditioner = 15
	maxTemperatureConditioner = 30
	minDepartments            = 1
	maxDepartments            = 1000
	minEmployees              = 1
	maxEmployees              = 1000
)

var (
	ErrInvalidOperator       = errors.New("invalid operator. Must be '>=' or '<='")
	ErrTemperatureOutOfRange = errors.New("temperature is out of range [15, 30]")
	ErrReadingDepartments    = errors.New("error reading departments count")
	ErrReadingEmployees      = errors.New("error reading employees count")
	ErrReadingInput          = errors.New("error reading operator and temperature")
	ErrDepartmentsOutOfRange = errors.New("departments is out of range [1, 1000]")
	ErrEmployeesOutOfRange   = errors.New("employees is out of range [1, 1000]")
)

type TemperatureRange struct {
	minTemp int
	maxTemp int
}

func NewTemperatureRange() *TemperatureRange {
	return &TemperatureRange{
		minTemp: minTemperatureConditioner,
		maxTemp: maxTemperatureConditioner,
	}
}

func (tr *TemperatureRange) Update(operator string, temperature int) error {
	if operator != ">=" && operator != "<=" {
		return ErrInvalidOperator
	}

	if temperature < minTemperatureConditioner || temperature > maxTemperatureConditioner {
		return ErrTemperatureOutOfRange
	}

	switch operator {
	case ">=":
		if temperature > tr.minTemp {
			tr.minTemp = temperature
		}
	case "<=":
		if temperature < tr.maxTemp {
			tr.maxTemp = temperature
		}
	}

	return nil
}

func (tr *TemperatureRange) GetResult() int {
	if tr.minTemp <= tr.maxTemp {
		return tr.minTemp
	}

	return -1
}

func processDepartment(employees int) error {
	tempRange := NewTemperatureRange()

	for range make([]struct{}, employees) {
		var (
			operator    string
			temperature int
		)

		_, err := fmt.Scanln(&operator, &temperature)
		if err != nil {
			return ErrReadingInput
		}

		err = tempRange.Update(operator, temperature)
		if err != nil {
			return err
		}

		fmt.Println(tempRange.GetResult())
	}

	return nil
}

func main() {
	var departments int

	_, err := fmt.Scanln(&departments)
	if err != nil {
		fmt.Println(ErrReadingDepartments)

		return
	}

	if departments < minDepartments || departments > maxDepartments {
		fmt.Println(ErrDepartmentsOutOfRange)

		return
	}

	for range make([]struct{}, departments) {
		var employees int

		_, err = fmt.Scanln(&employees)
		if err != nil {
			fmt.Println(ErrReadingEmployees)

			return
		}

		if employees < minEmployees || employees > maxEmployees {
			fmt.Println(ErrEmployeesOutOfRange)

			return
		}

		err = processDepartment(employees)
		if err != nil {
			fmt.Println(err)

			return
		}
	}
}
