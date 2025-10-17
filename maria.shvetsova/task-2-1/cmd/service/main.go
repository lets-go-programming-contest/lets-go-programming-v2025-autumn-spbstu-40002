package main

import (
	"errors"
	"fmt"

	"github.com/ummmsh/task-2-1/tempdata"
)

const (
	minDepEmpNum = 1
	maxDepEmpNum = 1000
	minTemp      = 15
	maxTemp      = 30
)

var (
	ErrInvalidDepartmentts = errors.New("invalid number of departments")
	ErrInvalidEmployees    = errors.New("invalid number of employees")
)

func main() {
	var numOfDepartments int

	_, err := fmt.Scan(&numOfDepartments)
	if err != nil || (numOfDepartments < minDepEmpNum || numOfDepartments > maxDepEmpNum) {
		fmt.Println(ErrInvalidDepartmentts)

		return
	}

	for range numOfDepartments {
		var numOfEmployees int

		_, err = fmt.Scan(&numOfEmployees)
		if err != nil || (numOfEmployees < minDepEmpNum || numOfEmployees > maxDepEmpNum) {
			fmt.Println(ErrInvalidEmployees)

			return
		}

		var sign string

		var temperature int

		optimalTemp := 15

		tempData, errNewTempData := tempdata.NewTempData(optimalTemp, maxTemp, minTemp)
		if errNewTempData != nil {
			fmt.Println(errNewTempData)

			return
		}

		for range numOfEmployees {
			_, err = fmt.Scan(&sign, &temperature)
			if err != nil {
				fmt.Println(tempdata.ErrInvalidSign)

				return
			}

			err = tempData.ChangeOptimalTemp(sign, temperature)
			if err != nil {
				fmt.Println(err)

				return
			}

			fmt.Println(tempData.GetOptimalTemp())
		}
	}
}
