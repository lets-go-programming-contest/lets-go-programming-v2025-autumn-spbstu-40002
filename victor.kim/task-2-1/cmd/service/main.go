package main

import (
	"errors"
	"fmt"
)

const (
	minTemperature                = 15
	maxTemperature                = 30
	expectedScanCountSingle       = 1
	expectedScanCountOperationVal = 2
)

var errFormat = errors.New("invalid temperature format")

func adjustTemperature(low int, high int) (int, int, error) {
	var operation string
	var newTemp int

	scanCount, scanErr := fmt.Scanln(&operation, &newTemp)
	if scanErr != nil {
		return 0, 0, errFormat
	}

	if scanCount != expectedScanCountOperationVal {
		return 0, 0, errFormat
	}

	if newTemp < minTemperature || newTemp > maxTemperature {
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

func main() {
	var (
		dep uint16
		emp uint16
	)

	scanCount, scanErr := fmt.Scanln(&dep)
	if scanErr != nil {
		fmt.Println("invalid number of departments")

		return
	}

	if scanCount != expectedScanCountSingle {
		fmt.Println("invalid number of departments")

		return
	}

	if dep == 0 || dep > 1000 {
		fmt.Println("invalid number of departments")

		return
	}

	for range make([]struct{}, dep) {
		scanCount, scanErr = fmt.Scanln(&emp)
		if scanErr != nil {
			fmt.Println("invalid number of employees")

			return
		}

		if scanCount != expectedScanCountSingle {
			fmt.Println("invalid number of employees")

			return
		}

		if emp == 0 || emp > 1000 {
			fmt.Println("invalid number of employees")

			return
		}

		lowTemp := minTemperature
		highTemp := maxTemperature

		for range make([]struct{}, emp) {
			var err error

			lowTemp, highTemp, err = adjustTemperature(lowTemp, highTemp)
			if err != nil {
				fmt.Println(err)

				return
			}

			fmt.Println(lowTemp)
		}
	}
}
