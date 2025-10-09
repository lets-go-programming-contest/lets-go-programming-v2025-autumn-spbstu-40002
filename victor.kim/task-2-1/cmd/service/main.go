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

func adjustTemperature(low int, high int) (int, int, error) {
	var (
		operation string
		newTemp   int
	)

	_, err := fmt.Scanln(&operation, &newTemp)
	if err != nil || newTemp < minTemperature || newTemp > maxTemperature {
		return 0, 0, errFormat
	}

	switch operation {
	case ">=":
		if high < newTemp {
			return -1, -1, nil
		}

		if low < newTemp {
			low = newTemp
		}
	case "<=":
		if low > newTemp {
			return -1, -1, nil
		}

		if high > newTemp {
			high = newTemp
		}
	default:
		return 0, 0, errFormat
	}

	return low, high, nil
}

func processDepartment(emp uint16) error {
	thermostat := struct {
		low      int
		high     int
		conflict bool
	}{
		low:      minTemperature,
		high:     maxTemperature,
	}

	for range make([]struct{}, emp) {
		var (
			operation string
			newTemp   int
		)

		_, err := fmt.Scanln(&operation, &newTemp)
		if err != nil || newTemp < minTemperature || newTemp > maxTemperature {
			return errFormat
		}

		switch operation {
		case ">=":
			if thermostat.high < newTemp {
				thermostat.conflict = true
			} else if thermostat.low < newTemp {
				thermostat.low = newTemp
			}
		case "<=":
			if thermostat.low > newTemp {
				thermostat.conflict = true
			} else if thermostat.high > newTemp {
				thermostat.high = newTemp
			}
		default:
			return errFormat
		}

		if thermostat.conflict {
			fmt.Println(-1)
		} else {
			fmt.Println(thermostat.low)
		}
	}

	return nil
}

func main() {
	var (
		dep uint16
		emp uint16
	)

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

		if procErr := processDepartment(emp); procErr != nil {
			fmt.Println(procErr)

			return
		}
	}
}
