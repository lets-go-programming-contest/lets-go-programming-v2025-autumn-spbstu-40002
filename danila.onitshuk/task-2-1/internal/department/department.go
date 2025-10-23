package department

import "fmt"

type Department struct {
	leftValueTemperatureRange  int
	rightValueTemperatureRange int
}

func NewDepartment() *Department {
	return &Department{
		leftValueTemperatureRange:  MinimumTemperature,
		rightValueTemperatureRange: MaximumTemperature,
	}
}

func (d *Department) ReduceRange(sign string, temperature int) error {
	switch sign {
	case ">=":
		if temperature > d.leftValueTemperatureRange {
			d.leftValueTemperatureRange = temperature
		}
	case "<=":
		if temperature < d.rightValueTemperatureRange {
			d.rightValueTemperatureRange = temperature
		}
	default:
		return ErrInvalidSign
	}

	if d.leftValueTemperatureRange <= d.rightValueTemperatureRange {
		fmt.Println(d.leftValueTemperatureRange)
	} else {
		fmt.Println(-1)
	}

	return nil
}
