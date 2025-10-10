package main

import "fmt"

const (
	minTemp = 15
	maxTemp = 30
)

func main() {
	var departmentCount int
	_, err := fmt.Scan(&departmentCount)
	if err != nil {
		return
	}

	for range departmentCount {
		var employeesCount int
		_, err = fmt.Scan(&employeesCount)
		if err != nil {
			return
		}

		low := minTemp
		high := maxTemp

		for range employeesCount {
			var operatorSign string
			var value int

			_, err = fmt.Scan(&operatorSign, &value)
			if err != nil {
				return
			}

			switch operatorSign {
			case ">=":
				if value > low {
					low = value
				}
			case "<=":
				if value < high {
					high = value
				}
			}

			if low > high {
				fmt.Println(-1)
			} else {
				fmt.Println(low)
			}
		}
	}
}
