package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Employee struct {
	minTemp int
	maxTemp int
}

const partsExpected = 2

func readInt(scanner *bufio.Scanner) (int, error) {
	scanner.Scan()

	return strconv.Atoi(strings.TrimSpace(scanner.Text()))
}

func readEmployees(scanner *bufio.Scanner, count int) []*Employee {
	employees := make([]*Employee, count)

	for empIndex := range employees {
		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != partsExpected {
			continue
		}

		value, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		employees[empIndex] = &Employee{minTemp: 0, maxTemp: 0}

		switch parts[0] {
		case ">=":
			employees[empIndex].minTemp = value
		case "<=":
			employees[empIndex].maxTemp = value
		default:
			continue
		}
	}

	return employees
}

func processDepartments(deptCount int, employees []*Employee) {
	for range make([]struct{}, deptCount) {
		commonMinTemp := 15
		commonMaxTemp := 30

		for _, emp := range employees {
			if emp == nil {
				continue
			}

			if emp.minTemp > commonMinTemp {
				commonMinTemp = emp.minTemp
			}

			if emp.maxTemp < commonMaxTemp {
				commonMaxTemp = emp.maxTemp
			}
		}

		if commonMinTemp <= commonMaxTemp {
			fmt.Println(commonMinTemp)
		} else {
			fmt.Println(-1)
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	departmentsCount, err := readInt(scanner)
	if err != nil || departmentsCount <= 0 {
		fmt.Fprintln(os.Stderr, "Ошибка чтения количества департаментов")

		return
	}

	employeesCount, err := readInt(scanner)
	if err != nil || employeesCount <= 0 {
		fmt.Fprintln(os.Stderr, "Ошибка чтения количества сотрудников")

		return
	}

	employees := readEmployees(scanner, employeesCount)

	processDepartments(departmentsCount, employees)
}
