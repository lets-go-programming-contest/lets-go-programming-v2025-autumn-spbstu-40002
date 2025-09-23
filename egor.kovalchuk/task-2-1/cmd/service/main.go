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
			var op string
			var temp int
			if _, err := fmt.Scan(&op, &temp); err != nil {
				return
			}

			if contradiction {
				fmt.Println(-1)

				// blank line before continue (nlreturn)
				continue
			}

			switch op {
			case ">=":
				if temp > low {
					low = temp
				}
			case "<=":
				if temp < high {
					high = temp
				}
			}

			if low > high {
				contradiction = true

				fmt.Println(-1)
				continue
			}

			fmt.Println(low)
		}
	}
}
