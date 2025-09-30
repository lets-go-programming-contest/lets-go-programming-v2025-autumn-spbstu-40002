package main

import (
	"fmt"
)

const (
	MIN = 15
	MAX = 30
)

func hateLint() {
	var amountOfDepartments int
	_, err := fmt.Scan(&amountOfDepartments)

	if err != nil || amountOfDepartments < 1 || amountOfDepartments > 1000 {
		fmt.Println("Incorrect amount of departments")

		return
	}

	var sign string

	newBorder := 0
	amountOfEmployees := 0

	for range amountOfDepartments {
		leftBorder := MIN
		rightBorder := MAX

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

			_, err = fmt.Scan(&sign, &newBorder)
			if err != nil {
				return
			} else if newBorder > MAX || newBorder < MIN {
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

func main() {
	hateLint()
}
