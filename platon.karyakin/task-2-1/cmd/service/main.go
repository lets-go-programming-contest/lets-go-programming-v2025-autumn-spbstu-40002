package main

import (
	"fmt"
	"log"
	"os"
)

const (
	initialMinTemp = 15
	initialMaxTemp = 30
)

var errFormat = fmt.Errorf("invalid input format or value")

func readOperationTemp() (string, int, error) {
	var operation string
	var temp int

	_, err := fmt.Scanln(&operation, &temp)
	if err != nil {
		return "", 0, errFormat
	}

	if operation != ">=" && operation != "<=" {
		return "", 0, errFormat
	}

	if temp < initialMinTemp || temp > initialMaxTemp {
		return "", 0, errFormat
	}

	return operation, temp, nil
}

func main() {
	var departmentsCount int

	if _, err := fmt.Fscan(os.Stdin, &departmentsCount); err != nil {
		return
	}

	for departmentIndex := 0; departmentIndex < departmentsCount; departmentIndex++ {
		var employeesCount int

		if _, err := fmt.Fscan(os.Stdin, &employeesCount); err != nil {
			return
		}

		currentMinTemp := initialMinTemp
		currentMaxTemp := initialMaxTemp

		for employeeIndex := 0; employeeIndex < employeesCount; employeeIndex++ {
			operation, temperature, err := readOperationTemp()
			if err != nil {

				log.Fatalf("error: %v", err)
			}

			if operation == ">=" {
				if temperature > currentMinTemp {
					currentMinTemp = temperature
				}
			} else if operation == "<=" {
				if temperature < currentMaxTemp {
					currentMaxTemp = temperature
				}
			}

			if currentMinTemp <= currentMaxTemp {
				fmt.Println(currentMinTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
