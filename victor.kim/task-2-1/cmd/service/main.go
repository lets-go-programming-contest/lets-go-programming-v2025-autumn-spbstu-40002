package main

import (
	"errors"
	"fmt"
)

const (
	minTemperature = 15
	maxTemperature = 30
)

var errFormat = errors.New("invalid temperature format")

type Thermostat struct {
	low      int
	high     int
	conflict bool
}

func NewThermostat() *Thermostat {
	return &Thermostat{
		low:      minTemperature,
		high:     maxTemperature,
		conflict: false,
	}
}

func (t *Thermostat) Update(operation string, value int) (bool, error) {
	if value < minTemperature || value > maxTemperature {
		return false, errFormat
	}

	if t.conflict {
		return true, nil
	}

	switch operation {
	case ">=":
		if value > t.high {
			t.conflict = true

			return true, nil
		}

		if value > t.low {
			t.low = value
		}
	case "<=":
		if value < t.low {
			t.conflict = true

			return true, nil
		}

		if value < t.high {
			t.high = value
		}
	default:
		return false, errFormat
	}

	return t.conflict, nil
}

func (t *Thermostat) Current() int {
	if t.conflict {
		return -1
	}

	return t.low
}

func readOperationAndValue() (string, int, error) {
	var operation string
	var value int

	_, err := fmt.Scanln(&operation, &value)
	if err != nil {
		return "", 0, errFormat
	}

	if value < minTemperature || value > maxTemperature {
		return "", 0, errFormat
	}

	return operation, value, nil
}

func processDepartment(employees uint16) error {
	thermostat := NewThermostat()

	for range make([]struct{}, employees) {
		operation, value, err := readOperationAndValue()
		if err != nil {
			return err
		}

		conflict, updErr := thermostat.Update(operation, value)
		if updErr != nil {
			return updErr
		}

		if conflict {
			fmt.Println(-1)
		} else {
			fmt.Println(thermostat.Current())
		}
	}

	return nil
}

func main() {
	var dep uint16
	var emp uint16

	_, err := fmt.Scanln(&dep)
	if err != nil || dep == 0 || dep > 1000 {
		fmt.Println("invalid number of departments")

		return
	}

	for range make([]struct{}, dep) {
		_, err = fmt.Scanln(&emp)
		if err != nil || emp == 0 || emp > 1000 {
			fmt.Println("invalid number of employees")

			return
		}

		if err := processDepartment(emp); err != nil {
			fmt.Println(err)

			return
		}
	}
}
