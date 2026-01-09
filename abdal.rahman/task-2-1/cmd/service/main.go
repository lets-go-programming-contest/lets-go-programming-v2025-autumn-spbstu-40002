package main

import (
	"fmt"
)

const (
	minTemperatureConditioner = 15
	maxTemperatureConditioner = 30
	minDepartments            = 1
	maxDepartments            = 1000
	minEmployees              = 1
	maxEmployees              = 1000
)

type TemperatureRange struct {
	minTemp int
	maxTemp int
}

func NewTemperatureRange() *TemperatureRange {
	return &TemperatureRange{
		minTemp: minTemperatureConditioner,
		maxTemp: maxTemperatureConditioner,
	}
}

func (tr *TemperatureRange) Update(operator string, temperature int) bool {
	if operator != ">=" && operator != "<=" {
		return false
	}

	if temperature < minTemperatureConditioner || temperature > maxTemperatureConditioner {
		return false
	}

	switch operator {
	case ">=":
		if temperature > tr.minTemp {
			tr.minTemp = temperature
		}
	case "<=":
		if temperature < tr.maxTemp {
			tr.maxTemp = temperature
		}
	}

	return true
}

func (tr *TemperatureRange) GetResult() int {
	if tr.minTemp <= tr.maxTemp {
		return tr.minTemp
	}

	return -1
}

func main() {
	var departments int
	_, err := fmt.Scan(&departments)
	if err != nil || departments < minDepartments || departments > maxDepartments {
		return
	}

	for i := 0; i < departments; i++ {
		var employees int
		_, err = fmt.Scan(&employees)
		if err != nil || employees < minEmployees || employees > maxEmployees {
			return
		}

		tempRange := NewTemperatureRange()

		for j := 0; j < employees; j++ {
			var operator string
			var temperature int

			_, err = fmt.Scan(&operator, &temperature)
			if err != nil {
				fmt.Println(-1)

				continue
			}

			ok := tempRange.Update(operator, temperature)
			if !ok {
				fmt.Println(-1)

			} else {
				fmt.Println(tempRange.GetResult())

			}
		}
	}
}
