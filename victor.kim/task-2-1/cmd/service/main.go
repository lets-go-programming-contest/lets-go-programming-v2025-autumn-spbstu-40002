package main

import (
	"errors"
	"fmt"
)

const (
	minTemperature = 15
	maxTemperature = 30
)

var (
	errFormat     = errors.New("invalid temperature format")
	errScanFailed = errors.New("scan failed")
)

type Thermostat struct {
	low      int
	high     int
	conflict bool
}

func NewThermostat() *Thermostat {
	return &Thermostat{
		low:  minTemperature,
		high: maxTemperature,
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

func (t *Thermostat) Current() (int, error) {
	if t.conflict {
		return -1, nil
	}

	return t.low, nil
}

func readInt() (int, error) {
	var value int
	scanCount, scanErr := fmt.Scan(&value)
	if scanErr != nil {
		return 0, errScanFailed
	}

	if scanCount != 1 {
		return 0, errScanFailed
	}

	return value, nil
}

func readOperationAndValue() (string, int, error) {
	var operation string
	var value int

	scanCount, scanErr := fmt.Scan(&operation, &value)
	if scanErr != nil {
		return "", 0, errScanFailed
	}

	if scanCount != 2 {
		return "", 0, errScanFailed
	}

	return operation, value, nil
}

func processDepartment(employees int) error {
	therm := NewThermostat()

	for range make([]struct{}, employees) {
		operation, value, readErr := readOperationAndValue()
		if readErr != nil {
			return readErr
		}

		conflict, updateErr := therm.Update(operation, value)
		if updateErr != nil {
			return updateErr
		}

		if conflict {
			fmt.Println(-1)
			continue
		}

		temp, curErr := therm.Current()
		if curErr != nil {
			return curErr
		}

		fmt.Println(temp)
	}

	return nil
}

func main() {
	departments, err := readInt()
	if err != nil {
		fmt.Println(err)
		return
	}

	if departments <= 0 || departments > 1000 {
		fmt.Println("invalid number of departments")
		return
	}

	for range make([]struct{}, departments) {
		employees, err := readInt()
		if err != nil {
			fmt.Println(err)
			return
		}

		if employees <= 0 || employees > 1000 {
			fmt.Println("invalid number of employees")
			return
		}

		if procErr := processDepartment(employees); procErr != nil {
			fmt.Println(procErr)
			return
		}
	}
}
