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

func readInt(scanner *bufio.Scanner) (int, error) {
	scanner.Scan()
	return strconv.Atoi(strings.TrimSpace(scanner.Text()))
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

	employees := make([]*Employee, employeesCount)
	for i := range employees {
		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		value, convErr := strconv.Atoi(parts[1])
		if convErr != nil {
			continue
		}

		employees[i] = &Employee{}

		switch parts[0] {
		case ">=":
			employees[i].minTemp = value
		case "<=":
			employees[i].maxTemp = value
		default:
			continue
		}
	}

	for range make([]struct{}, departmentsCount) {
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
