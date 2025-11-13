package main

import (
	"errors"
	"fmt"
)

const (
	MinTemp = 15
	MaxTemp = 30
	MinData = 1
	MaxData = 1000
)

var (
	ErrReadingDepartments    = errors.New("error reading departments count")
	ErrReadingEmployees      = errors.New("error reading employees count")
	ErrReadingInput          = errors.New("error reading operator and temperature")
	ErrDepartmentsOutOfRange = errors.New("departments is out of range [1, 1000]")
	ErrEmployeesOutOfRange   = errors.New("employees is out of range [1, 1000]")
)

type TempRange struct {
	Min int
	Max int
}

func isValid(data int) bool {
	return data >= MinData && data <= MaxData
}

func findBestTemp(tr TempRange, operator string, temp int) TempRange {
	switch operator {
	case ">=":
		if temp > tr.Min {
			tr.Min = temp
		}
	case "<=":
		if temp < tr.Max {
			tr.Max = temp
		}
	}
	return tr
}

func main() {
	var otdels int

	_, err := fmt.Scan(&otdels)
	if err != nil {
		fmt.Println(ErrReadingDepartments)
		return
	}

	if !isValid(otdels) {
		fmt.Println(ErrDepartmentsOutOfRange)
		return
	}

	for range otdels {
		tr := TempRange{Min: MinTemp, Max: MaxTemp}

		var workers int

		_, err := fmt.Scan(&workers)
		if err != nil {
			fmt.Println(ErrReadingEmployees)
			return
		}
		if !isValid(workers) {
			fmt.Println(ErrEmployeesOutOfRange)
			return
		}

		results := make([]int, workers)

		for i := range workers {
			var operator string
			var temp int

			_, err = fmt.Scan(&operator)
			if err != nil {
				fmt.Println(ErrReadingInput)
				return
			}

			_, err = fmt.Scan(&temp)
			if err != nil {
				fmt.Println(ErrReadingInput)
				return
			}

			tr = findBestTemp(tr, operator, temp)

			if tr.Min <= tr.Max {
				results[i] = tr.Min
			} else {
				results[i] = -1
			}
		}

		for _, result := range results {
			fmt.Println(result)
		}
	}
}
