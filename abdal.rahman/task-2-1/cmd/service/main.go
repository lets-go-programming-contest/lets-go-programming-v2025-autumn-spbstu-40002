package main

import (
	"errors"
	"fmt"
)

const (
	TempMin = 15
	TempMax = 30

	MinDepartments = 1
	MaxDepartments = 1000
)

var (
	ErrBadOperator     = errors.New("bad operator")
	ErrTempOutOfRange  = errors.New("temperature out of range")
	ErrReadDepartments = errors.New("failed to read departments")
	ErrReadEmployees   = errors.New("failed to read employees")
)

type DeptTemp struct {
	Min int
	Max int
}

func (d *DeptTemp) Update(op string, temp int) error {
	switch op {
	case ">=":
		if temp > d.Min {
			d.Min = temp
		}
	case "<=":
		if temp < d.Max {
			d.Max = temp
		}
	default:
		return ErrBadOperator
	}

	return nil
}

func (d *DeptTemp) Current() (int, error) {
	if d.Min > d.Max {
		return -1, ErrTempOutOfRange
	}

	return d.Min, nil
}

func main() {
	var departments int

	if _, err := fmt.Scan(&departments); err != nil ||
		departments < MinDepartments || departments > MaxDepartments {
		fmt.Println("Error:", ErrReadDepartments)

		return
	}

	for range departments {
		var employees int

		if _, err := fmt.Scan(&employees); err != nil {
			fmt.Println("Error:", ErrReadEmployees)

			return
		}

		dept := DeptTemp{
			Min: TempMin,
			Max: TempMax,
		}

		for range employees {
			var op string
			var temp int

			if _, err := fmt.Scan(&op, &temp); err != nil {
				fmt.Println("Error: invalid input")

			}

			if err := dept.Update(op, temp); err != nil {
				fmt.Println("Error:", err)

				return
			}

			value, err := dept.Current()
			if err != nil {
				fmt.Println(-1)
			} else {
				fmt.Println(value)
			}
		}
	}
}
