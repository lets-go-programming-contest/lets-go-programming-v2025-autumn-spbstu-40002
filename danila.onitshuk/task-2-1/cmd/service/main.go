package main

import (
	"fmt"

	department "danila.onitshuk/task-2-1/internal/department"
)

func main() {
	var departmentsCount int

	_, err := fmt.Scan(&departmentsCount)
	if err != nil {
		fmt.Println(department.ErrInvalidTypeData)

		return
	} else if departmentsCount < department.MinimumCountDepartments ||
		departmentsCount > department.MaximumCountDepartments {
		fmt.Println(department.ErrCountDepartments)

		return
	}

	for range departmentsCount {
		var employeesCount int

		_, err := fmt.Scan(&employeesCount)
		if err != nil {
			fmt.Println(department.ErrInvalidTypeData)

			return
		} else if employeesCount < department.MinimumCountEmployees || employeesCount > department.MaximumCountEmployees {
			fmt.Println(department.ErrCountEmployees)

			return
		}

		currentDepartment := department.NewDepartment()

		for range employeesCount {
			var (
				sign        string
				temperature int
			)

			_, err := fmt.Scan(&sign, &temperature)
			if err != nil {
				fmt.Println(department.ErrInvalidTypeData)

				return
			}

			err = currentDepartment.ReduceRange(sign, temperature)
			if err != nil {
				fmt.Println(err)

				return
			}
		}
	}
}
