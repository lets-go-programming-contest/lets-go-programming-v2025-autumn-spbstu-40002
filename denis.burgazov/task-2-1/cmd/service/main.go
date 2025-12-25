package main

import (
	"fmt"
)

func parse(a any, errText string) bool {
	_, err := fmt.Scanln(a)
	if err != nil {
		fmt.Println(errText)

		var errBuf string

		_, _ = fmt.Scanln(&errBuf)
	}

	return err == nil
}

func main() {
	var (
		numberOfDept int
		numberOfEmp  int
	)

	if !parse(&numberOfDept, "Invalid number of department") {
		return
	}

	for range numberOfDept {
		if !parse(&numberOfEmp, "Invalid number of employee") {
			return
		}

		var (
			minTemp  = 15
			maxTemp  = 30
			sign     string
			currTemp int
		)

		for range numberOfEmp {
			_, _ = fmt.Scan(&sign)

			if !parse(&currTemp, "Invalid temperature value") {
				return
			}

			if minTemp == -1 {
				fmt.Println(-1)

				continue
			}

			switch sign {
			case ">=":
				if currTemp <= maxTemp {
					minTemp = max(currTemp, minTemp)
				} else {
					minTemp = -1
				}
			case "<=":
				if currTemp >= minTemp {
					maxTemp = min(currTemp, maxTemp)
				} else {
					minTemp = -1
				}
			default:
				fmt.Println("Invalid input")

				return
			}

			fmt.Println(minTemp)
		}
	}
}
