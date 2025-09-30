package main

import (
	"fmt"
)

const (
	TEMPERATURE_MIN = 15
	TEMPERATURE_MAX = 30
	DEP_EMP_MAX     = 1000
	DEP_EMP_MIN     = 1
)

func loopForSpecificDepartment(leftBorder int, rightBorder int, amountOfEmployees int) error {
	var newBorder int
	var sign string

	for range amountOfEmployees {
		if leftBorder == -1 {
			fmt.Println(leftBorder)

			continue
		}

		_, err := fmt.Scan(&sign)
		if err != nil {
			return fmt.Errorf("incorrect sign")
		}

		_, err = fmt.Scan(&newBorder)
		if err != nil || newBorder > TEMPERATURE_MAX || newBorder < TEMPERATURE_MIN {
			return fmt.Errorf("incorrect border")
		}

		switch sign {
		case "<=":
			if newBorder < leftBorder {
				leftBorder = -1
			} else if newBorder < rightBorder {
				rightBorder = newBorder
			}
		case ">=":
			if newBorder > rightBorder {
				leftBorder = -1
			} else if newBorder > leftBorder {
				leftBorder = newBorder
			}
		default:
			return fmt.Errorf("incorrect sign")
		}

		fmt.Println(leftBorder)
	}
	return nil
}

func main() {
	var amountOfDepartments int
	_, err := fmt.Scan(&amountOfDepartments)

	if err != nil || amountOfDepartments < DEP_EMP_MIN || amountOfDepartments > DEP_EMP_MAX {
		fmt.Println("Error: incorrect amount of departments")

		return
	}

	amountOfEmployees := 0

	for range amountOfDepartments {
		_, err = fmt.Scan(&amountOfEmployees)
		if err != nil || amountOfEmployees < DEP_EMP_MIN || amountOfEmployees > DEP_EMP_MAX {
			fmt.Println("Error: incorrect amount of employees")

			return
		}

		err = loopForSpecificDepartment(TEMPERATURE_MIN, TEMPERATURE_MAX, amountOfEmployees)
		if err != nil {
			fmt.Println("Error:", err)

			return
		}
	}
}
