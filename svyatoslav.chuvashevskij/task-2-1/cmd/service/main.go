package main

import (
	"errors"
	"fmt"
)

const (
	MaxTemperature = 30
	MinTemperature = 15
)

var ErrInvalidOption = errors.New("invalid option")

type Condition struct {
	minTemperature int
	maxTemperature int
}

func (condition *Condition) SetTemperature(option string, temperature int) error {
	switch option {
	case ">=":
		if condition.minTemperature < temperature {
			condition.minTemperature = temperature
		}
	case "<=":
		if condition.maxTemperature > temperature {
			condition.maxTemperature = temperature
		}
	default:
		return ErrInvalidOption
	}

	return nil
}

func (condition *Condition) getTemperature() int {
	if condition.maxTemperature < condition.minTemperature {
		return -1
	}

	return condition.minTemperature
}

func main() {
	numberOfDepartments := 0
	numberOfEmployees := 0

	_, err := fmt.Scan(&numberOfDepartments)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	for range numberOfDepartments {
		_, err = fmt.Scan(&numberOfEmployees)
		if err != nil {
			fmt.Println("Error:", err)

			return
		}

		condition := Condition{minTemperature: MinTemperature, maxTemperature: MaxTemperature}

		for range numberOfEmployees {
			var option string

			var temperature int

			_, err = fmt.Scan(&option, &temperature)
			if err != nil {
				fmt.Println("Error:", err)

				return
			}

			err = condition.SetTemperature(option, temperature)
			if err != nil {
				fmt.Println("Error:", err)

				return
			}

			fmt.Println(condition.getTemperature())
		}
	}
}
