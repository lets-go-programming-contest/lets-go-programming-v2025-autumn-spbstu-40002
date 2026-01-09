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

func (t *Temperature) Adjust(operator string, value int) error {
	if value < MinTemp || value > MaxTemp {
		return ErrTempOutOfRange
	}

	switch operator {
	case ">=":
		if value > t.Min {
			t.Min = value
		}
	case "<=":
		if value < t.Max {
			t.Max = value
		}
	default:
		return ErrInvalidSign
	}
	return nil
}

func main() {
	var numDepts int
	_, err := fmt.Scan(&numDepts)
	if err != nil {
		fmt.Println(ErrReadingInput)
		return
	}

	if numDepts < MinRange || numDepts > MaxRange {
		fmt.Println(ErrDeptOutOfRange)
		return
	}

	for d := 0; d < numDepts; d++ {
		var numEmps int
		_, err := fmt.Scan(&numEmps)
		if err != nil {
			fmt.Println(ErrReadingInput)
			return
		}

		if numEmps < MinRange || numEmps > MaxRange {
			fmt.Println(ErrEmpOutOfRange)
			return
		}

		tempRange, err := NewTemperature(MaxTemp, MinTemp)
		if err != nil {
			fmt.Println(err)
			return
		}

		for e := 0; e < numEmps; e++ {
			var op string
			var val int

			_, err := fmt.Scan(&op, &val)
			if err != nil {
				fmt.Println(ErrReadingInput)
				return
			}

			if err := tempRange.Adjust(op, val); err != nil {
				fmt.Println(-1)
			} else {
				fmt.Println(tempRange.Optimal())
			}
		}
	}
}
