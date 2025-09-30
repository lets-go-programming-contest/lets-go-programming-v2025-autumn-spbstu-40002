package main

import "fmt"

const (
	defaultLow  = 15
	defaultHigh = 30
)

func main() {
	var departments int
	if _, err := fmt.Scan(&departments); err != nil || departments < 0 {
		return
	}

	depts := make([]struct{}, departments)
	for range depts {
		if ok := processDepartment(); !ok {
			return
		}
	}
}

func processDepartment() bool {
	low, high := defaultLow, defaultHigh

	var employeeCount int
	if _, err := fmt.Scan(&employeeCount); err != nil || employeeCount < 0 {
		return false
	}

	contradiction := false

	var operator string
	var temperature int

	emps := make([]struct{}, employeeCount)
	for range emps {

		if _, err := fmt.Scan(&operator, &temperature); err != nil {
			return false
		}

		if contradiction {
			fmt.Println(-1)

			continue
		}

		switch operator {
		case ">=":
			if temperature > low {
				low = temperature
			}
		case "<=":
			if temperature < high {
				high = temperature
			}
		default:
			contradiction = true
			fmt.Println(-1)

			continue
		}

		if low > high {
			contradiction = true
			fmt.Println(-1)

			continue
		}

		fmt.Println(low)
	}

	return true
}
