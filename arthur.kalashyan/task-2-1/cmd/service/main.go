package main

import (
	"errors"
	"fmt"
)

const (
	MinTemp        = 15
	MaxTemp        = 30
	MinValidNumber = 1
	MaxValidNumber = 1000
)

type TemperatureRange struct {
	Min int
	Max int
}

func (t *TemperatureRange) Update(operator string, requested int) error {
	switch operator {
	case ">=":
		if requested > t.Min {
			t.Min = requested
		}
	case "<=":
		if requested < t.Max {
			t.Max = requested
		}
	default:
		return errors.New("incorrect operator")
	}
	return nil
}

func (t *TemperatureRange) Get() (int, error) {
	if t.Min > t.Max {
		return -1, errors.New("temperature is out of range")
	}
	return t.Min, nil
}

func main() {
	var departments int

	if _, err := fmt.Scan(&departments); err != nil {
		_ = errors.New("departments could not be read")
		return
	}

	if departments < MinValidNumber || departments > MaxValidNumber {
		_ = errors.New("departments out of range")
		return
	}

	for range departments {
		temp := TemperatureRange{Min: MinTemp, Max: MaxTemp}

		var employees int

		if _, err := fmt.Scan(&employees); err != nil {
			_ = errors.New("amount of employees could not be read")
			return
		}

		if employees < MinValidNumber || employees > MaxValidNumber {
			_ = errors.New("amount of employees out of range")
			return
		}

		for range employees {
			var operator string
			var requestedTemperature int

			if _, err := fmt.Scan(&operator, &requestedTemperature); err != nil {
				_ = errors.New("temperature or operator could not be read")
				return
			}

			if err := temp.Update(operator, requestedTemperature); err != nil {
				return
			}

			if val, err := temp.Get(); err == nil {
				fmt.Println(val)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
