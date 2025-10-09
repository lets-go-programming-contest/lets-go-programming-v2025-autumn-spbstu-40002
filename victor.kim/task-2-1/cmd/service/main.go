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

func (t *Thermostat) Update(operation string, newTemp int) (bool, error) {
	if newTemp < minTemperature || newTemp > maxTemperature {
		return false, errFormat
	}

	if t.conflict {
		return true, nil
	}

	switch operation {
	case ">=":
		if t.high < newTemp {
			t.conflict = true

			return true, nil
		}

		if t.low < newTemp {
			t.low = newTemp
		}
	case "<=":
		if t.low > newTemp {
			t.conflict = true

			return true, nil
		}

		if t.high > newTemp {
			t.high = newTemp
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

		// === minimal change: include all struct fields in composite literal ===
		t := &Thermostat{
			low:      minTemperature,
			high:     maxTemperature,
			conflict: false,
		}
		// ====================================================================

		for range make([]struct{}, emp) {
			var operation string
			var newTemp int

			_, err := fmt.Scanln(&operation, &newTemp)
			if err != nil || newTemp < minTemperature || newTemp > maxTemperature {
				fmt.Println(errFormat)

				return
			}

			conflict, updErr := t.Update(operation, newTemp)
			if updErr != nil {
				fmt.Println(updErr)

				return
			}

			if conflict {
				fmt.Println(-1)
			} else {
				fmt.Println(t.Current())
			}
		}
	}
}
