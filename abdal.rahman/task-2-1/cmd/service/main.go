package main

import (
	"errors"
	"fmt"
)

const (
	MinTemp  = 15
	MaxTemp  = 30
	MinRange = 1
	MaxRange = 1000
)

type Temperature struct {
	Min int
	Max int
}

var (
	ErrInvalidSign    = errors.New("invalid sign")
	ErrTempOutOfRange = errors.New("temperature out of range")
	ErrDeptOutOfRange = errors.New("departments out of range")
	ErrEmpOutOfRange  = errors.New("employees out of range")
	ErrReadingInput   = errors.New("error reading input")
)

func NewTemperature(maxTemp, minTemp int) (*Temperature, error) {
	if minTemp > maxTemp {
		return nil, ErrTempOutOfRange
	}

	return &Temperature{
		Max: maxTemp,
		Min: minTemp,
	}, nil
}

func (t *Temperature) Optimal() int {
	if t.Min > t.Max {
		return -1
	}

	return t.Min
}

func (t *Temperature) Adjust(operatorStr string, temperatureVal int) error {
	if temperatureVal < MinTemp || temperatureVal > MaxTemp {
		return ErrTempOutOfRange
	}

	switch operatorStr {
	case ">=":
		if temperatureVal > t.Min {
			t.Min = temperatureVal
		}
	case "<=":
		if temperatureVal < t.Max {
			t.Max = temperatureVal
		}
	default:
		return ErrInvalidSign
	}

	return nil
}

func main() {
	var totalDepts int
	_, err := fmt.Scan(&totalDepts)
	if err != nil {
		fmt.Println(ErrReadingInput)
		return
	}

	if totalDepts < MinRange || totalDepts > MaxRange {
		fmt.Println(ErrDeptOutOfRange)
		return
	}

	for dept := 0; dept < totalDepts; dept++ {
		var totalEmps int
		_, err := fmt.Scan(&totalEmps)
		if err != nil {
			fmt.Println(ErrReadingInput)
			return
		}

		if totalEmps < MinRange || totalEmps > MaxRange {
			fmt.Println(ErrEmpOutOfRange)
			return
		}

		tempRange, err := NewTemperature(MaxTemp, MinTemp)
		if err != nil {
			fmt.Println(err)
			return
		}

		for emp := 0; emp < totalEmps; emp++ {
			var operatorStr string
			var temperatureVal int

			_, err := fmt.Scan(&operatorStr, &temperatureVal)
			if err != nil {
				fmt.Println(ErrReadingInput)
				return
			}

			err = tempRange.Adjust(operatorStr, temperatureVal)
			if err != nil {
				fmt.Println(-1)
			} else {
				fmt.Println(tempRange.Optimal())
			}
		}
	}
}
