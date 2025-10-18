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

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	departmentsCount, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || departmentsCount <= 0 {
		panic("Ошибка чтения количества департаментов")
	}

	scanner.Scan()
	employeesCount, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || employeesCount <= 0 {
		panic("Ошибка чтения количества сотрудников")
	}

	employees := make([]*Employee, employeesCount)

	for empIndex := 0; empIndex < employeesCount; empIndex++ {
		scanner.Scan()
		line := scanner.Text()
		parts := strings.Fields(line)

		switch parts[0] {
		case ">=":
			value, _ := strconv.Atoi(parts[1])
			employees[empIndex].minTemp = value
		case "<=":
			value, _ := strconv.Atoi(parts[1])
			employees[empIndex].maxTemp = value
		default:
			continue
		}
	}

	for deptIndex := 0; deptIndex < departmentsCount; deptIndex++ {
		commonMinTemp := 15
		commonMaxTemp := 30

		for _, employee := range employees {
			if commonMinTemp < employee.minTemp {
				commonMinTemp = employee.minTemp
			}
			if commonMaxTemp > employee.maxTemp {
				commonMaxTemp = employee.maxTemp
			}
		}

		if commonMinTemp <= commonMaxTemp {
			fmt.Printf("%d\n", commonMinTemp)
		} else {
			fmt.Printf("-1\n")
		}

		if deptIndex+1 < departmentsCount {
			fmt.Println("")
		}
	}
}
