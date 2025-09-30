package main

import (
	"errors"
	"fmt"
)

const (
	TemperatureMin = 15
	TemperatureMax = 30
	DepEmpMin      = 1
	DepEmpMax      = 1000
)

var (
	ErrorIncorrectSign        = errors.New("incorrect sign")
	ErrorIncorrectBorder      = errors.New("incorrect border")
	ErrorIncorrectDepartments = errors.New("incorrect amount of departments")
	ErrorIncorrectEmployees   = errors.New("incorrect amount of employees")
)

func loopForSpecificDepartment(leftBorder int, rightBorder int, amountOfEmployees int) error {
	var newBorder int

	var sign string

	for range amountOfEmployees {

		_, err := fmt.Scan(&sign)
		if err != nil {
			return ErrorIncorrectSign
		}

		_, err = fmt.Scan(&newBorder)
		if err != nil || newBorder > TemperatureMax || newBorder < TemperatureMin {
			return ErrorIncorrectBorder
		}

		if leftBorder == -1 {
			fmt.Println(leftBorder)

			continue
		}

		switch sign {
		case "<=":
			if newBorder < rightBorder {
				rightBorder = newBorder
			}
		case ">=":
			if newBorder > leftBorder {
				leftBorder = newBorder
			}
		default:
			return ErrorIncorrectSign
		}

		if leftBorder > rightBorder {
			leftBorder = -1
		}

		fmt.Println(leftBorder)
	}
	return nil
}

func main() {
	var amountOfDepartments int
	_, err := fmt.Scan(&amountOfDepartments)

	if err != nil || amountOfDepartments < DepEmpMin || amountOfDepartments > DepEmpMax {
		fmt.Println("Error:", ErrorIncorrectDepartments)

		return
	}

	amountOfEmployees := 0

	for range amountOfDepartments {
		_, err = fmt.Scan(&amountOfEmployees)
		if err != nil || amountOfEmployees < DepEmpMin || amountOfEmployees > DepEmpMax {
			fmt.Println("Error:", ErrorIncorrectEmployees)

			return
		}

		err = loopForSpecificDepartment(TemperatureMin, TemperatureMax, amountOfEmployees)
		if err != nil {
			fmt.Println("Error:", err)

			return
		}
	}
}
