package main

import (
	"fmt"
)

const (
	MIN = 15
	MAX = 30
)

func main() {
	var amountOfDepartments int
	_, err := fmt.Scan(&amountOfDepartments)

	if err != nil || amountOfDepartments < 1 || amountOfDepartments > 1000 {
		fmt.Println("Incorrect amount of departments")

		return
	}

	var sign string
	leftBorder := 0
	rightBorder := 0
	newBorder := 0
	amountOfEmployees := 0

	for range amountOfDepartments {
		leftBorder = MIN
		rightBorder = MAX

		_, err = fmt.Scan(&amountOfEmployees)
		if err != nil || amountOfEmployees < 1 || amountOfEmployees > 1000 {
			fmt.Println("Incorrect amount of employees")

			return
		}

		for range amountOfEmployees {
			if leftBorder == -1 {
				fmt.Println(leftBorder)

				continue
			}

			_, err = fmt.Scan(&sign)
			if err != nil {
				return
			}

			_, err = fmt.Scan(&newBorder)
			if err != nil || newBorder > 30 || newBorder < 15 {
				fmt.Println("Incorrect border")

				return
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
				fmt.Println("Incorrect sign")

				return
			}

			fmt.Println(leftBorder)
		}
	}
}
