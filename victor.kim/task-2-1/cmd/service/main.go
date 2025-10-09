package main

import (
	"errors"
	"fmt"
)

const (
	minTemperature                = 15
	maxTemperature                = 30
	expectedScanCountOperationVal = 2
	expectedScanCountSingle       = 1
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

func readOperationAndValue() (string, int, error) {
	var (
		operation string
		value     int
	)

	scanCount, scanErr := fmt.Scanln(&operation, &value)
	if scanErr != nil {
		return "", 0, errFormat
	}

	if scanCount != expectedScanCountOperationVal {
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
	var (
		dep uint16
		emp uint16
	)

	scanCount, scanErr := fmt.Scanln(&dep)
	if scanErr != nil {
		fmt.Println("invalid number of departments")

		return
	}

	if scanCount != expectedScanCountSingle {
		fmt.Println("invalid number of departments")

		return
	}

	if dep == 0 || dep > 1000 {
		fmt.Println("invalid number of departments")

		return
	}

	for range make([]struct{}, dep) {
		scanCount, scanErr = fmt.Scanln(&emp)
		if scanErr != nil {
			fmt.Println("invalid number of employees")

			return
		}

		if scanCount != expectedScanCountSingle {
			fmt.Println("invalid number of employees")

			return
		}

		if emp == 0 || emp > 1000 {
			fmt.Println("invalid number of employees")

			return
		}

		if procErr := processDepartment(emp); procErr != nil {
			fmt.Println(procErr)

			return
		}
	}
}
