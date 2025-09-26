package main

import "fmt"

func main() {
	var departments int
	if _, err := fmt.Scan(&departments); err != nil {
		return
	}

	for range make([]struct{}, departments) {
		low := 15
		high := 30

		var employeeCount int
		if _, err := fmt.Scan(&employeeCount); err != nil {
			return
		}

		contradiction := false

		for range make([]struct{}, employeeCount) {
			operator := ""
			temperature := 0

			_, err := fmt.Scan(&operator, &temperature)
			if err != nil {
				return
			}

			if contradiction {
				fmt.Println(-1)

				// пустая строка перед continue (nlreturn)
				continue
			}

			switch operator {
			case ">=":
				if temperature > low {
					low = temperature
				}
			case "<=":
				if temperature < high {
					high = temperature
				}
			}

			if low > high {
				contradiction = true

				fmt.Println(-1)

				// пустая строка перед continue (nlreturn)
				continue
			}

			fmt.Println(low)
		}
	}
}
