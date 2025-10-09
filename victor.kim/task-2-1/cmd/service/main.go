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

func (th *Thermostat) Update(operation string, value int) (bool, error) {
	if value < minTemperature || value > maxTemperature {
		return false, errFormat
	}

	if th.conflict {
		return true, nil
	}

	switch operation {
	case ">=":
		if value > th.high {
			th.conflict = true
			return true, nil
		}

		if value > th.low {
			th.low = value
		}
	case "<=":
		if value < th.low {
			th.conflict = true
			return true, nil
		}

		if value < th.high {
			th.high = value
		}
	default:
		return false, errFormat
	}

	return th.conflict, nil
}

func (th *Thermostat) Current() int {
	if th.conflict {
		return -1
	}

	return th.low
}

func processDepartment(employees int) error {
	thermostat := NewThermostat()

	for range make([]struct{}, employees) {
		var (
			operation string
			newTemp   int
		)

		scanCount, scanErr := fmt.Scanln(&operation, &newTemp)
		if scanErr != nil {
			return errFormat
		}

		if scanCount != 2 {
			return errFormat
		}

		if newTemp < minTemperature || newTemp > maxTemperature {
			return errFormat
		}

		conflict, updErr := thermostat.Update(operation, newTemp)
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
	var departments int

	scanCount, scanErr := fmt.Scanln(&departments)
	if scanErr != nil {
		fmt.Println("invalid number of departments")

		return
	}

	if scanCount != 1 {
		fmt.Println("invalid number of departments")

		return
	}

	if departments <= 0 || departments > 1000 {
		fmt.Println("invalid number of departments")

		return
	}

	for range make([]struct{}, departments) {
		var employees int

		scanCount, scanErr = fmt.Scanln(&employees)
		if scanErr != nil {
			fmt.Println("invalid number of employees")

			return
		}

		if scanCount != 1 {
			fmt.Println("invalid number of employees")

			return
		}

		if employees <= 0 || employees > 1000 {
			fmt.Println("invalid number of employees")

			return
		}

		if err := processDepartment(employees); err != nil {
			fmt.Println(err)

			return
		}
	}
}
