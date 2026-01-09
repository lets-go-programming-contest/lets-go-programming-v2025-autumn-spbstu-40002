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

type TempRange struct {
	Min int
	Max int
}

var (
	ErrInvalidSign     = errors.New("invalid sign")
	ErrTempOutOfBounds = errors.New("temperature out of range")
	ErrDeptOutOfBounds = errors.New("departments out of range")
	ErrEmpOutOfBounds  = errors.New("employees out of range")
	ErrReadingInput    = errors.New("error reading input")
)

func newTempRange(maxTemp, minTemp int) (*TempRange, error) {
	if minTemp > maxTemp {
		return nil, ErrTempOutOfBounds
	}

	return &TempRange{
		Max: maxTemp,
		Min: minTemp,
	}, nil
}

func (tr *TempRange) Optimal() int {
	if tr.Min > tr.Max {
		return -1
	}
	return tr.Min
}

func (tr *TempRange) Adjust(operator string, value int) error {
	if value < MinTemp || value > MaxTemp {
		return ErrTempOutOfBounds
	}

	switch operator {
	case ">=":
		if value > tr.Min {
			tr.Min = value
		}
	case "<=":
		if value < tr.Max {
			tr.Max = value
		}
	default:
		return ErrInvalidSign
	}
	return nil
}

func main() {
	var totalDepts int
	if _, err := fmt.Scan(&totalDepts); err != nil {
		fmt.Println(ErrReadingInput)
		return
	}

	if totalDepts < MinRange || totalDepts > MaxRange {
		fmt.Println(ErrDeptOutOfBounds)
		return
	}

	for dept := 0; dept < totalDepts; dept++ {
		var totalEmps int
		if _, err := fmt.Scan(&totalEmps); err != nil {
			fmt.Println(ErrReadingInput)
			return
		}

		if totalEmps < MinRange || totalEmps > MaxRange {
			fmt.Println(ErrEmpOutOfBounds)
			return
		}

		tempRange, err := newTempRange(MaxTemp, MinTemp)
		if err != nil {
			fmt.Println(err)
			return
		}

		for emp := 0; emp < totalEmps; emp++ {
			var operator string
			var value int

			if _, err := fmt.Scan(&operator, &value); err != nil {
				fmt.Println(ErrReadingInput)
				return
			}

			if err := tempRange.Adjust(operator, value); err != nil {
				fmt.Println(-1)
			} else {
				fmt.Println(tempRange.Optimal())
			}
		}
	}
}
