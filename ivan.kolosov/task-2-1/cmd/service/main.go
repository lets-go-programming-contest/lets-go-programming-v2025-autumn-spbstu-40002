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
	ErrIncorrectSign        = errors.New("incorrect sign")
	ErrIncorrectBorder      = errors.New("incorrect border")
	ErrIncorrectDepartments = errors.New("incorrect amount of departments")
	ErrIncorrectEmployees   = errors.New("incorrect amount of employees")
)

func getSpecificEmployeeOpinion(leftBorder int, rightBorder int, newBorder int, sign string) (int, int, error) {
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
		return leftBorder, rightBorder, ErrIncorrectSign
	}

	if leftBorder > rightBorder {
		leftBorder = -1
	}

	return leftBorder, rightBorder, nil
}

func loopForSpecificDepartment(leftBorder int, rightBorder int, amountOfEmployees int) error {
	var newBorder int

	var sign string

	for range amountOfEmployees {
		_, err := fmt.Scan(&sign)
		if err != nil {
			return ErrIncorrectSign
		}

		_, err = fmt.Scan(&newBorder)
		if err != nil || newBorder > TemperatureMax || newBorder < TemperatureMin {
			return ErrIncorrectBorder
		}

		if leftBorder == -1 {
			fmt.Println(leftBorder)

			continue
		}

		leftBorder, rightBorder, err = getSpecificEmployeeOpinion(leftBorder, rightBorder, newBorder, sign)
		if err != nil {
			return err
		}

		fmt.Println(leftBorder)
	}

	return nil
}

func main() {
	var amountOfDepartments int
	_, err := fmt.Scan(&amountOfDepartments)

	if err != nil || amountOfDepartments < DepEmpMin || amountOfDepartments > DepEmpMax {
		fmt.Println("Error:", ErrIncorrectDepartments)

		return
	}

	amountOfEmployees := 0

	for range amountOfDepartments {
		_, err = fmt.Scan(&amountOfEmployees)
		if err != nil || amountOfEmployees < DepEmpMin || amountOfEmployees > DepEmpMax {
			fmt.Println("Error:", ErrIncorrectEmployees)

			return
		}

		err = loopForSpecificDepartment(TemperatureMin, TemperatureMax, amountOfEmployees)
		if err != nil {
			fmt.Println("Error:", err)

			return
		}
	}
}
