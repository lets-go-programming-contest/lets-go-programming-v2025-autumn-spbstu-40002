package main

import (
	"errors"
	"fmt"
)

const (
	minDeptNum   = 1
	maxDeptNum   = 1000
	minEmployees = 1
	maxEmployees = 1000
	minTemp      = 15
	maxTemp      = 30
)

var ErrUnknownOperation = errors.New("unknown operation")

type TempController struct {
	optimalTemp int
	minimumTemp int
	maximumTemp int
}

func (tc *TempController) adjustTemp(operation string, val int) error {
	switch operation {
	case ">=":
		if tc.minimumTemp < val {
			tc.minimumTemp = val
		}

		if tc.minimumTemp < minTemp {
			tc.minimumTemp = minTemp
		}

		if tc.optimalTemp < tc.minimumTemp {
			tc.optimalTemp = tc.minimumTemp
		}
	case "<=":
		if tc.maximumTemp > val {
			tc.maximumTemp = val
		}

		if tc.maximumTemp > maxTemp {
			tc.maximumTemp = maxTemp
		}

		if tc.optimalTemp > tc.maximumTemp {
			tc.optimalTemp = tc.maximumTemp
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnknownOperation, operation)
	}

	if tc.maximumTemp < tc.minimumTemp {
		tc.optimalTemp = -1
	}

	return nil
}

func (tc *TempController) currentTemp() int {
	return tc.optimalTemp
}

func NewTempController(low int, high int) *TempController {
	return &TempController{
		optimalTemp: low,
		minimumTemp: low,
		maximumTemp: high,
	}
}

func main() {
	var deptCount int

	_, err := fmt.Scanln(&deptCount)
	if err != nil || deptCount > maxDeptNum || deptCount < minDeptNum {
		fmt.Println("Invalid department count")

		return
	}

	deptIndex := 0
	for deptIndex < deptCount {
		var employeeCount int

		_, err = fmt.Scanln(&employeeCount)
		if err != nil || employeeCount > maxEmployees || employeeCount < minEmployees {
			fmt.Println("Invalid employee count")

			return
		}

		controller := NewTempController(minTemp, maxTemp)

		empIndex := 0
		for empIndex < employeeCount {
			var operation string
			var value int

			_, err = fmt.Scanln(&operation, &value)
			if err != nil {
				fmt.Println("Invalid input")

				return
			}

			if err := controller.adjustTemp(operation, value); err != nil {
				fmt.Println("Error:", err)

				return
			}

			currentTemp := controller.currentTemp()
			fmt.Println(currentTemp)

			empIndex++
		}

		deptIndex++
	}
}
