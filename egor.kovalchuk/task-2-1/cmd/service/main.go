package main

import (
	"errors"
	"fmt"
)

var errOperation = errors.New("invalid operation")

func adjustTemperature(low int, high int, temp int, op string) (int, int, error) {
	if low == -1 && high == -1 {
		return low, high, nil
	}

	switch op {
	case ">=":
		if temp > high {
			return -1, -1, nil
		}
		if temp > low {
			low = temp
		}
	case "<=":
		if temp < low {
			return -1, -1, nil
		}
		if temp < high {
			high = temp
		}
	default:
		return low, high, errOperation
	}

	return low, high, nil
}

func processDepartment(employeeCount int, minTemp int, maxTemp int) {
	low := minTemp
	high := maxTemp

	var op string
	var temp int

	emps := make([]struct{}, employeeCount)
	for range emps {

		_, err := fmt.Scan(&op, &temp)
		if err != nil || temp < minTemp || temp > maxTemp {
			fmt.Println(-1)

			return
		}

		low, high, err = adjustTemperature(low, high, temp, op)
		if err != nil {
			fmt.Println(-1)

			return
		}

		fmt.Println(low)
	}
}

func main() {
	const (
		minTemp = 15
		maxTemp = 30
	)

	var departmentCount, employeeCount int

	_, err := fmt.Scan(&departmentCount)
	if err != nil || departmentCount < 1 {
		return
	}

	depts := make([]struct{}, departmentCount)
	for range depts {

		_, err := fmt.Scan(&employeeCount)
		if err != nil || employeeCount < 0 {
			return
		}

		processDepartment(employeeCount, minTemp, maxTemp)
	}
}
